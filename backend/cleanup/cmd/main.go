package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
	"github.com/thyamix/sumcrowds/cleanup/internal/cleanup"
)

func main() {
	getEnv()
	db.InitDB()
	cleanup.Clean()
}

func getEnv() {
	if os.Getenv("APP_DEPLOY") != "docker" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found (set APP_DEPLOY to 'docker' if deployed with docker)")
		}
	}
}
