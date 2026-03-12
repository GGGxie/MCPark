package main

import (
	"log"
	"net/http"

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

	cfg := config.Load()
	database := db.InitDB(cfg.DatabasePath)
	defer database.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL, "http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	dedup := handler.NewDedup()
	authHandler := handler.NewAuthHandler(database, cfg)
	serviceHandler := handler.NewServiceHandler(database, dedup)
	adminHandler := handler.NewAdminHandler(database)
	articleHandler := handler.NewArticleHandler(database, dedup)

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

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/google", authHandler.GoogleLogin)
			auth.GET("/google/callback", authHandler.GoogleCallback)
			auth.POST("/google/token", authHandler.GoogleTokenLogin)
			auth.GET("/me", middleware.AuthRequired(cfg.JWTSecret), authHandler.GetMe)
		}

		services := api.Group("/services")
		{
			services.GET("", serviceHandler.ListServices)
			services.GET("/featured", serviceHandler.GetFeaturedServices)
			services.GET("/categories", serviceHandler.GetCategories)
			services.GET("/:id", serviceHandler.GetService)
			services.POST("", middleware.AuthRequired(cfg.JWTSecret), serviceHandler.CreateService)
			services.POST("/:id/like", middleware.AuthRequired(cfg.JWTSecret), serviceHandler.LikeService)
			services.POST("/:id/install", serviceHandler.InstallService)
			services.GET("/:id/comments", serviceHandler.GetComments)
			services.POST("/:id/comments", middleware.AuthRequired(cfg.JWTSecret), serviceHandler.CreateComment)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AdminRequired(cfg.JWTSecret))
		{
			admin.GET("/services", adminHandler.ListAllServices)
			admin.POST("/services/:id/approve", adminHandler.ApproveService)
			admin.POST("/services/:id/reject", adminHandler.RejectService)
			admin.DELETE("/services/:id", adminHandler.DeleteService)
		}

		articles := api.Group("/articles")
		{
			articles.GET("", articleHandler.ListArticles)
			articles.GET("/:id", articleHandler.GetArticle)
		}
	}

	log.Printf("MCPark server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
