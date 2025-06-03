package main

import (
	"log"
	// "net/http" // No longer directly needed for ListenAndServe with Gin

	"fipe_project/internal/database"
	"fipe_project/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.ConnectMongoDB(); err != nil {
		log.Fatalf("Erro ao conectar no MongoDB: %v", err)
	}

	// Initialize Gin router
	router := gin.Default() // gin.Default() comes with Logger and Recovery middleware

	// Setup routes using the modified SetupRoutes function
	routes.SetupRoutes(router) // SetupRoutes now expects a *gin.Engine

	log.Println("Servidor Gin rodando na porta 8080")
	// router.Run() is the Gin equivalent of http.ListenAndServe
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar servidor Gin: %v", err)
	}
}


