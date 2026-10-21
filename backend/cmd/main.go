package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"MCPark/internal/config"
	"MCPark/internal/db"
	"MCPark/internal/handler"
	"MCPark/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Must run after .env is loaded so GIN_MODE from .env is honoured.
	gin.SetMode(resolveGinMode(os.Getenv))

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	database := db.InitDB(cfg.DatabasePath)
	defer database.Close()

	r, err := newEngine(cfg, database)
	if err != nil {
		log.Fatal("Failed to build router: ", err)
	}

	addr := listenAddress(cfg)
	log.Printf("MCPark server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// resolveGinMode defaults to release mode unless GIN_MODE is set explicitly.
func resolveGinMode(getenv func(string) string) string {
	if m := getenv("GIN_MODE"); m != "" {
		return m
	}
	return gin.ReleaseMode
}

// listenAddress joins LISTEN_ADDR (default 127.0.0.1) and PORT.
func listenAddress(cfg *config.Config) string {
	return net.JoinHostPort(cfg.ListenAddr, cfg.Port)
}

func newEngine(cfg *config.Config, database *sql.DB) (*gin.Engine, error) {
	r := gin.Default()

	// Only a reverse proxy on the same host may set client address headers.
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, err
	}

	r.Use(middleware.MaxBodyBytes(middleware.DefaultMaxBodyBytes))
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
	return r, nil
}
