package handler

import (
	"database/sql"

	"MCPark/internal/config"
	"MCPark/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts the JSON API under /api on r.
func RegisterRoutes(r *gin.Engine, database *sql.DB, cfg *config.Config) {
	dedup := NewDedup()
	authHandler := NewAuthHandler(database, cfg)
	serviceHandler := NewServiceHandler(database, dedup)
	adminHandler := NewAdminHandler(database)
	articleHandler := NewArticleHandler(database, dedup)

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
}
