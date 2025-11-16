package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB holds the database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB creates a new database connection pool
func NewDB(databaseURL string) (*DB, error) {
	// Parse connection string and create pool config
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 25                  // Maximum number of connections
	config.MinConns = 5                   // Minimum number of connections
	config.MaxConnLifetime = time.Hour    // Max lifetime of a connection
	config.MaxConnIdleTime = 30 * time.Minute // Max idle time before closing
	config.HealthCheckPeriod = time.Minute    // Health check interval

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Println("Database connection closed")
	}
}

// HealthCheck checks if database is accessible
func (db *DB) HealthCheck(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
