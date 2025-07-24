package http

import (
	"net/http"

	api_handler_v1 "github.com/thyamix/sumcrowds/backend/counter/internal/net/http/handlers/v1"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/http/middleware"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/websockets"
)

func getRoutes(wsServer *websockets.Server) *http.ServeMux {
	routes := http.NewServeMux()

	//WS Routes

	routes.HandleFunc("/ws/{festivalCode}", func(w http.ResponseWriter, r *http.Request) { websockets.HandleCounter(wsServer, w, r) })

	//V1 Routes

	routes.HandleFunc("GET /api/v1/auth/validateaccess", api_handler_v1.ValidateAccess)
	routes.HandleFunc("GET /api/v1/auth/refreshaccess", api_handler_v1.RefreshAccess)
	routes.HandleFunc("GET /api/v1/auth/initaccess", api_handler_v1.InitAccess)

	routes.Handle("POST /api/v1/create/festival", middleware.RequireAccess(http.HandlerFunc(api_handler_v1.CreateFestival)))

	// /festival behind auth
	routes.Handle("GET /api/v1/festival/{festivalCode}/access", middleware.RequireAccess(http.HandlerFunc(api_handler_v1.CheckAccess)))

	routes.Handle("GET /api/v1/festival/{festivalCode}/exists", middleware.RequireAccess(http.HandlerFunc(api_handler_v1.CheckExists)))
	routes.Handle("POST /api/v1/festival/{festivalCode}/getaccess", middleware.RequireAccess(http.HandlerFunc(api_handler_v1.GetAccess)))

	routes.Handle("POST /api/v1/festival/{festivalCode}/inc", middleware.RequireAccess(http.HandlerFunc(api_handler_v1.Inc)))
	routes.Handle("POST /api/v1/festival/{festivalCode}/dec", middleware.RequireAccess(http.HandlerFunc(api_handler_v1.Dec)))

	routes.Handle("GET /api/v1/festival/{festivalCode}/totalandgauge", middleware.RequireAccess(http.HandlerFunc(api_handler_v1.GetTotalAndGauge)))

	//admin

	routes.Handle("GET /api/v1/festival/{festivalCode}/admin/access", middleware.RequireAccess(http.HandlerFunc(api_handler_v1.CheckAdminAccess)))

	routes.Handle("GET /api/v1/festival/{festivalCode}/admin/getarchivedevents", middleware.RequireAccess(middleware.RequiresAdminPin(http.HandlerFunc(api_handler_v1.GetArchivedEvents))))

	routes.Handle("POST /api/v1/festival/{festivalCode}/admin/setgauge", middleware.RequireAccess(middleware.RequiresAdminPin(http.HandlerFunc(api_handler_v1.SetGauge))))

	routes.Handle("GET /api/v1/festival/{festivalCode}/admin/archivecurrentevent", middleware.RequireAccess(middleware.RequiresAdminPin(http.HandlerFunc(api_handler_v1.ArchiveCurrentEvent))))

	routes.Handle("GET /api/v1/festival/{festivalCode}/admin/download/archivedcsv/{eventId}", middleware.RequireAccess(middleware.RequiresAdminPin(http.HandlerFunc(api_handler_v1.GetArchivedCSV))))

	routes.Handle("GET /api/v1/festival/{festivalCode}/admin/download/activecsv", middleware.RequireAccess(middleware.RequiresAdminPin(http.HandlerFunc(api_handler_v1.GetActiveCSV))))

	return routes
}
