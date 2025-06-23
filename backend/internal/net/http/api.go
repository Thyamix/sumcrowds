package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/thyamix/festival-counter/internal/net/websockets"
)

func StartAPI() {
	fmt.Println("Started StartAPI")
	wsServer := websockets.NewHub()
	router := getRoutes(wsServer)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ORIGIN")},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}).Handler(router)

	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: handler,
	}

	fmt.Println(os.Getenv("PORT"))

	fmt.Println("Server is starting...")

	go wsServer.Run()

	log.Fatal(server.ListenAndServe())
}
