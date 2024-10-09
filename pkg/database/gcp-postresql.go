package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/oauth2/google"
)

type GCPPostgresStrategy struct {
	db        *sql.DB
	gcpConfig GCPConfig
}

// NewGCPPostgresStrategy creates a new GCP PostgreSQL strategy
func NewGCPPostgresStrategy(gcpConfig GCPConfig) *GCPPostgresStrategy {
	return &GCPPostgresStrategy{
		gcpConfig: gcpConfig,
	}
}

// Connect implements the connection logic for GCP PostgreSQL
func (s *GCPPostgresStrategy) Connect(config Config) (*sql.DB, error) {
	ctx := context.Background()

	// if s.gcpConfig.UseIAMAuth {
	// 	password, err = s.generateIAMAuthToken(ctx)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to generate IAM auth token: %w", err)
	// 	}
	// 	config.Password = password
	// }

	dsn := s.BuildDSN(config)
	instanceConnectionName := s.BuildInstanceName()

	// db, err := sql.Open("postgres", dsn)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to open database: %w", err)
	// }
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	var opts []cloudsqlconn.Option
	if s.gcpConfig.UsePrivateIP {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	// Use the Cloud SQL connector to handle connecting to the instance.
	// This approach does *NOT* require the Cloud SQL proxy.
	connConfig.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return db, nil

	// Configure connection pool
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	// Verify connection
	if err := s.Ping(ctx, db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	s.db = db
	return db, nil
}

// BuildDSN creates the connection string for GCP PostgreSQL
func (s *GCPPostgresStrategy) BuildDSN(config Config) string {
	// var host string
	// if s.gcpConfig.UsePrivateIP {
	// 	host = s.gcpConfig.InstanceName
	// } else {
	// 	host = fmt.Sprintf("/cloudsql/%s", instanceConnectionName)
	// }

	return fmt.Sprintf("user=%s password=%s database=%s",
		config.User,
		config.Password,
		config.Database,
	)
}

func (s *GCPPostgresStrategy) BuildInstanceName() string {
	instanceConnectionName := fmt.Sprintf("%s:%s:%s",
		s.gcpConfig.ProjectID,
		s.gcpConfig.Region,
		s.gcpConfig.InstanceName,
	)

	return instanceConnectionName

}

// generateIAMAuthToken generates a token for IAM authentication
func (s *GCPPostgresStrategy) generateIAMAuthToken(ctx context.Context) (string, error) {
	var creds *google.Credentials
	var err error

	if s.gcpConfig.CredentialsFile != "" {
		credBytes, err := os.ReadFile(s.gcpConfig.CredentialsFile)
		if err != nil {
			return "", fmt.Errorf("failed to read credentials file: %w", err)
		}

		creds, err = google.CredentialsFromJSON(ctx, credBytes, "https://www.googleapis.com/auth/sqlservice.admin")
	} else {
		creds, err = google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/sqlservice.admin")
	}

	if err != nil {
		return "", fmt.Errorf("failed to get credentials: %w", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token.AccessToken, nil
}

// Ping verifies the database connection
func (s *GCPPostgresStrategy) Ping(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

// Close closes the database connection
func (s *GCPPostgresStrategy) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Transaction represents a database transaction
type Transaction struct {
	tx *sql.Tx
}

// ExecuteInTransaction executes the given function within a transaction
func (db *Database) ExecuteInTransaction(ctx context.Context, fn func(*Transaction) error) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(&Transaction{tx: tx}); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
