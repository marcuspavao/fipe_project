package routes

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	projecthandlers "fipe_project/internal/handlers"
)

func SetupRoutes() http.Handler {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/tabelas", projecthandlers.GetTabelasReferencia).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/marcas", projecthandlers.GetMarcas).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/modelos/{marca}", projecthandlers.GetModelos).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/veiculos", projecthandlers.GetVeiculos).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/dashboard", projecthandlers.GetDashboardMarcas).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/0km", projecthandlers.GetVeiculosNovos).Methods("GET", "OPTIONS")

	staticFileServer := http.FileServer(http.Dir("./frontend/"))
	router.PathPrefix("/").Handler(staticFileServer)

	allowedOriginsVal := []string{
		"http://localhost:5173",
		"http://192.168.3.31:5173",
	}

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
	)

	return corsHandler(router)
}
