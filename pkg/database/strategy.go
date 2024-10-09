package database

import (
	"context"
	"database/sql"
)

// DBStrategy defines the interface for database operations
type DBStrategy interface {
	Connect(config Config) (*sql.DB, error)
	BuildDSN(config Config) string
	Ping(ctx context.Context, db *sql.DB) error
	Close() error
}
