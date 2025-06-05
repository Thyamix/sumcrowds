package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/thyamix/festival-counter/internal/net/websockets"
)

func StartAPI() {
	wsServer := websockets.NewHub()
	router := getRoutes(wsServer)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://testing.sumcrowds.com"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}).Handler(router)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: handler,
	}

	fmt.Println("Server is starting...")

	go wsServer.Run()

	log.Fatal(server.ListenAndServe())
}
