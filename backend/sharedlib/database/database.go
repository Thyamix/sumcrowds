package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

//go:embed init.sql
var InitSQL string

var DB *sql.DB

func InitDB() {
	fmt.Println("Initialising Database")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Error: DATABASE_URL environment variable is not set.")
	}

	var err error

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
