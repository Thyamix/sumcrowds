package main

import (
	"log"

	"github.com/thyamix/sumcrowds/backend/counter/internal/net/http"
	"github.com/thyamix/sumcrowds/backend/sharedlib/config"
	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
)

func main() {
	// Load configuration
	env := config.GetEnv()
	cfg, err := config.Load(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting with %s configuration", env)

	// Initialize database with config
	db.InitDBWithConfig(cfg)

	// Start API server
	http.StartAPI(cfg)
}
