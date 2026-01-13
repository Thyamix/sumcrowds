package main

import (
	"log"

	"github.com/thyamix/sumcrowds/backend/sharedlib/config"
	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
	"github.com/thyamix/sumcrowds/cleanup/internal/cleanup"
)

func main() {
	// Load configuration
	env := config.GetEnv()
	cfg, err := config.Load(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting cleanup with %s configuration", env)

	// Initialize database with config
	db.InitDBWithConfig(cfg)

	// Run cleanup
	cleanup.Clean()
}
