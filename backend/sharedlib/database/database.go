package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/thyamix/sumcrowds/backend/sharedlib/config"
	"github.com/thyamix/sumcrowds/backend/sharedlib/database/sqlcdb"
)

//go:embed init.sql
var InitSQL string

var DB *sql.DB
var Pool *pgxpool.Pool
var Queries *sqlcdb.Queries

// InitDBWithConfig initializes the database using the config system
func InitDBWithConfig(cfg *config.Config) {
	fmt.Println("Initialising Database")

	connStr := cfg.GetDatabaseURL()
	initWithConnStr(connStr)
}

// InitDB initializes the database using legacy environment variable
// Deprecated: Use InitDBWithConfig instead
func InitDB() {
	fmt.Println("Initialising Database")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Error: DATABASE_URL environment variable is not set.")
	}
	initWithConnStr(connStr)
}

func initWithConnStr(connStr string) {

	var err error

	// Initialize legacy database/sql connection for backward compatibility
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DB.PingContext(pingCtx)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres database (ping failed): %v", err)
	}

	// Initialize pgx pool for SQLC
	poolCtx, poolCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer poolCancel()

	Pool, err = pgxpool.New(poolCtx, connStr)
	if err != nil {
		log.Fatalf("Failed to create pgx pool: %v", err)
	}

	err = Pool.Ping(poolCtx)
	if err != nil {
		log.Fatalf("Failed to ping pgx pool: %v", err)
	}

	// Initialize SQLC queries
	Queries = sqlcdb.New(Pool)

	fmt.Println("Successfully connected to database")

	err = initTables()
	if err != nil {
		log.Fatalf("Failed to initialise database with error: %v", err)
	}

	fmt.Println("Database succesfully initialised")
	fmt.Println("Database started")
}

func initTables() error {
	_, err := DB.Exec(InitSQL)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			return fmt.Errorf("PostgreSQL error during initTables: Code=%s, Detail=%s, %w", pgErr.Code.Name(), pgErr.Detail, err)
		}
		return fmt.Errorf("failed to execute init_pg.sql: %w", err)
	}

	return nil
}

// CloseDB closes both database connections (legacy sql.DB and pgx pool)
// This should be called during graceful shutdown
func CloseDB() {
	if Pool != nil {
		Pool.Close()
		fmt.Println("pgx pool closed")
	}
	if DB != nil {
		if err := DB.Close(); err != nil {
			fmt.Printf("Error closing legacy DB connection: %v\n", err)
		} else {
			fmt.Println("Legacy DB connection closed")
		}
	}
}
