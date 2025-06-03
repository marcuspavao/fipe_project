package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fipe_project/internal/handlers" // Adjusted import path
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes group
	api := router.Group("/api")
	{
		api.GET("/tabelas", handlers.GetTabelasReferencia) // Assuming GetTabelasReferencia is the correct name
		api.GET("/marcas", handlers.GetMarcas)             // Assuming GetMarcas is the correct name
		// The original Gorilla Mux route was /modelos/{marca}
		// The original JS call was /api/modelos/${marcaVal}?tabela=${tabelaVal}
		// Assuming GetModelos handler can take marca from query param or path, and tabela from query.
		// For Gin, if marca is a path param: api.GET("/modelos/:marca", handlers.GetModelos)
		// If marca is a query param: api.GET("/modelos", handlers.GetModelos)
		// Let's stick to a query param approach for marca to align with the JS call, assuming handler supports it.
		api.GET("/modelos", handlers.GetModelos) // Handler needs to get 'marca' and 'tabela' from query
		api.GET("/veiculos", handlers.GetVeiculos)
		api.GET("/dashboard", handlers.GetDashboardMarcas) // Assuming GetDashboardMarcas is correct
	}

	// Static file serving for Vue app
	distDir := "./frontend/dist"

	// Serve static files (js, css, images, etc.) from assets directory
	router.StaticFS("/assets", http.Dir(filepath.Join(distDir, "assets")))
	
	// Serve specific files from the root of distDir
	// Favicon might be /favicon.ico or /assets/favicon.ico depending on Vite build.
	// Check actual 'dist' structure after Vue build.
	// Typically Vite places generated assets in 'assets' folder, and index.html references them.
	// Public files can be in root of dist.
	router.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))
	// router.StaticFile("/images/favicon.png", filepath.Join(distDir, "images/favicon.png")) // This was the old path. Vue build places it in root or assets.

	// For any other route, serve index.html - Vue Router will handle the rest
	router.NoRoute(func(c *gin.Context) {
		// Check if it's an API call that wasn't matched
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"message": "API endpoint not found"})
			return
		}

		// Serve index.html for Vue app
		indexPath := filepath.Join(distDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			// This can happen if 'frontend/dist' is not populated correctly.
			c.String(http.StatusInternalServerError, "index.html not found at %s", indexPath)
			return
		}
		c.File(indexPath)
	})
}
