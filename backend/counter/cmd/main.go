package main

import (
	"github.com/joho/godotenv"
	"github.com/thyamix/sumcrowds/backend/counter/internal/database"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/http"
	"log"
)

func main() {
	getEnv()
	database.InitDB()
	http.StartAPI()
}

func getEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (set APP_DEPLOY to 'docker' if deployed with docker)")
	}
}
