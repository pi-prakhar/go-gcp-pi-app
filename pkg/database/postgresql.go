package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresStrategy struct {
	db *sql.DB
}

func NewPostgresStrategy() *PostgresStrategy {
	return &PostgresStrategy{}
}

func (s *PostgresStrategy) Connect(config Config) (*sql.DB, error) {
	dsn := s.BuildDSN(config)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	s.db = db
	s.configureConnection(config)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Ping(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *PostgresStrategy) BuildDSN(config Config) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

	// Add additional options from config
	for key, value := range config.Options {
		dsn += fmt.Sprintf(" %s=%s", key, value)
	}

	return dsn
}

func (s *PostgresStrategy) Ping(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

func (s *PostgresStrategy) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *PostgresStrategy) configureConnection(config Config) {
	s.db.SetMaxOpenConns(config.MaxOpenConns)
	s.db.SetMaxIdleConns(config.MaxIdleConns)
	s.db.SetConnMaxLifetime(config.MaxLifetime)
}
