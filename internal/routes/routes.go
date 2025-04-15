package routes

import (
	"net/http"
	"fipe_project/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/tabelas", handlers.GetTabelasReferencia).Methods("GET")
	router.HandleFunc("/api/marcas", handlers.GetMarcas).Methods("GET")
	router.HandleFunc("/api/modelos/{marca}", handlers.GetModelos).Methods("GET")
	router.HandleFunc("/api/veiculos", handlers.GetVeiculos).Methods("GET")
	router.HandleFunc("/api/dashboard", handlers.GetDashboardMarcas).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))

	return router
}
