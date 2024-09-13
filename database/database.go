package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDatabaseConnection() {
	if DB == nil {
		var db_url = os.Getenv("DB_URL")

		var err error
		DB, err = GetDbPoolSqlx(db_url)

		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	} else {
		log.Println("Using existing database connection...")
	}

}

// Open connection with Sqlx
func GetDbPoolSqlx(uri string) (*sqlx.DB, error) {

	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	db.DB.SetMaxOpenConns(30)           // The default is 0 (unlimited)
	db.DB.SetMaxIdleConns(10)           // defaultMaxIdleConns = 2
	db.DB.SetConnMaxLifetime(time.Hour) // 0, connections are reused forever.

	return db, nil
}

// Open connection with pgxpool
func GetDbPoolSqlxPgx(uri string) (*sqlx.DB, error) {
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 30
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	// Test the connection
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Println("Failed to acquire connection from pool:", err)
		return nil, err
	}

	// To be safe
	defer conn.Release()

	// Get a connection
	db := stdlib.OpenDBFromPool(pool)
	if db == nil {
		return nil, fmt.Errorf("failed to open db")
	}

	return sqlx.NewDb(db, "postgres"), nil
}
