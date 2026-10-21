package main

import (
	"log"
	"net/http"

	"MCPark/internal/config"
	"MCPark/internal/db"
	"MCPark/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Invalid configuration: ", err)
	}
	database := db.InitDB(cfg.DatabasePath)
	defer database.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL, "http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Serve frontend static files
	r.Static("/assets", "./static/assets")
	r.StaticFile("/vite.svg", "./static/vite.svg")
	r.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.File("./static/index.html")
	})

	handler.RegisterRoutes(r, database, cfg)

	log.Printf("MCPark server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
