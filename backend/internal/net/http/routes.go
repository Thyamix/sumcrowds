package http

import (
	"net/http"

	api_handler_v1 "github.com/thyamix/festival-counter/internal/net/http/handlers/v1"
	"github.com/thyamix/festival-counter/internal/net/websockets"
)

func getRoutes(wsServer *websockets.Server) *http.ServeMux {
	routes := http.NewServeMux()

	//WS Routes

	routes.HandleFunc("/ws/{festivalCode}", func(w http.ResponseWriter, r *http.Request) { websockets.HandleCounter(wsServer, w, r) })

	//V1 Routes

	routes.HandleFunc("GET /api/v1/auth/validateaccess", api_handler_v1.ValidateAccess)
	routes.HandleFunc("GET /api/v1/auth/refreshaccess", api_handler_v1.RefreshAccess)
	routes.HandleFunc("GET /api/v1/auth/initaccess", api_handler_v1.InitAccess)

	routes.HandleFunc("POST /api/v1/create/festival", api_handler_v1.CreateFestival)

	// /festival behind auth

	routes.HandleFunc("GET /api/v1/festival/{festivalCode}/exists", api_handler_v1.CheckExists)
	routes.HandleFunc("GET /api/v1/festival/{festivalCode}/access", api_handler_v1.CheckAccess)
	routes.HandleFunc("POST /api/v1/festival/{festivalCode}/getAccess", api_handler_v1.GetAccess)

	routes.HandleFunc("GET /api/v1/festival/{festivalCode}/totalandgauge", api_handler_v1.GetTotalAndGauge)
	routes.HandleFunc("GET /api/v1/festival/{festivalCode}/getarchivedevents", api_handler_v1.GetArchivedEvents)

	routes.HandleFunc("POST /api/v1/festival/{festivalCode}/setgauge", api_handler_v1.SetGauge)

	routes.HandleFunc("GET /api/v1/festival/{festivalCode}/archivecurrentevent", api_handler_v1.ArchiveCurrentEvent)

	routes.HandleFunc("GET /api/v1/festival/{festivalCode}/download/archivedcsv/{eventId}", api_handler_v1.GetArchivedCSV)

	routes.HandleFunc("GET /api/v1/festival/{festivalCode}/download/activecsv", api_handler_v1.GetActiveCSV)

	return routes
}
