package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/websockets"
	"github.com/thyamix/sumcrowds/backend/sharedlib/config"
)

func StartAPI(cfg *config.Config) {
	fmt.Println("Starting API")
	wsServer := websockets.NewHub()
	router := getRoutes(wsServer)

	// Build allowed origins map for fast lookup
	allowedOriginsMap := make(map[string]bool)
	for _, origin := range cfg.CORS.AllowedOrigins {
		allowedOriginsMap[origin] = true
	}

	handler := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			// Allow requests with no Origin header (mobile apps)
			if origin == "" {
				return true
			}
			// Allow configured origins
			return allowedOriginsMap[origin]
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}).Handler(router)

	server := http.Server{
		Addr:    cfg.GetServerAddr(),
		Handler: handler,
	}

	go wsServer.Run()

	fmt.Printf("Server is listening to origins: %v\n", cfg.CORS.AllowedOrigins)
	fmt.Printf("Server is running on %s...\n", cfg.GetServerAddr())
	log.Fatal(server.ListenAndServe())
}
