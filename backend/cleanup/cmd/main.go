package main

import (
	"log"

	"github.com/thyamix/sumcrowds/backend/cleanup/internal/cleanup"
	"github.com/thyamix/sumcrowds/backend/sharedlib/config"
	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
)

func main() {
	log.Println("Starting...")
	// Load configuration
	env := config.GetEnv()
	log.Printf("Starting cleanup with %s configuration", env)
	cfg, err := config.Load(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database with config
	db.InitDBWithConfig(cfg)

	// Run cleanup
	cleanup.Clean()
}
