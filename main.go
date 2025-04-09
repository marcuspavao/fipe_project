package main

import (
	"log"
	"net/http"
    
	"fipe_project/internal/database"
	"fipe_project/internal/routes"
)

func main() {

	if err := database.ConnectMongoDB(); err != nil {
		log.Fatalf("Erro ao conectar no MongoDB: %v", err)
	}

	router := routes.SetupRoutes()

	log.Println("Servidor rodando na porta 8080 \n")
	log.Fatal(http.ListenAndServe(":8080", router))
}


