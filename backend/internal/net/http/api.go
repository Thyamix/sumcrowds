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
	fmt.Println("Starting API")
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

	go wsServer.Run()

	fmt.Printf("Server is listening to origin %v \n", os.Getenv("ORIGIN"))
	fmt.Printf("Server is running on :%v...\n", os.Getenv("PORT"))
	log.Fatal(server.ListenAndServe())
}
