package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/http"
	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
)

func main() {
	getEnv()
	db.InitDB()
	http.StartAPI()
}

func getEnv() {
	if os.Getenv("APP_DEPLOY") != "docker" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found (set APP_DEPLOY to 'docker' if deployed with docker)")
		}
	}
}
