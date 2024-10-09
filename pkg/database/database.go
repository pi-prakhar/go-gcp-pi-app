package database

import (
	"database/sql"
	"fmt"
	"sync"
)

// Database represents a database instance
type Database struct {
	db       *sql.DB
	config   Config
	strategy DBStrategy
}

var (
	instances = make(map[string]*Database)
	mu        sync.RWMutex
)

// NewDatabase creates or returns an existing database instance
func NewDatabase(name string, config Config, strategy DBStrategy) (*Database, error) {
	mu.Lock()
	defer mu.Unlock()

	if db, exists := instances[name]; exists {
		return db, nil
	}

	// strategy, err := getStrategy(config.Type)
	// if err != nil {
	// 	return nil, err
	// }

	db, err := strategy.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	instance := &Database{
		db:       db,
		config:   config,
		strategy: strategy,
	}
	instances[name] = instance

	return instance, nil
}

func (d *Database) Close() {
	d.strategy.Close()
}

// func getStrategy(dbType string) (DBStrategy, error) {
// 	switch dbType {
// 	// case "mysql":
// 	// 	return &MySQLStrategy{}, nil
// 	case "postgres":
// 		return &PostgresStrategy{}, nil
// 	default:
// 		return nil, fmt.Errorf("unsupported database type: %s", dbType)
// 	}
// }
