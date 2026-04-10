package main

import (
	"assignment6/handlers"
	_ "assignment6/docs" // Import generated docs
	"assignment6/middleware"
	"assignment6/models"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Assignment 6 API
// @version         1.0
// @description     API for Assignment 6 with JWT, RBAC and Rate Limiting.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type Bearer followed by a space and then your token.

func main() {
	r := gin.Default()

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initial User B (as admin) for testing TASK 2
	models.Mu.Lock()
	models.Users["admin@example.com"] = &models.User{
		ID:         999,
		Email:      "admin@example.com",
		Password:   "admin123",
		Role:       models.RoleAdmin,
		IsVerified: true,
	}
	models.Mu.Unlock()

	// Global Middleware
	r.Use(middleware.OptionalJWTAuthMiddleware())
	r.Use(middleware.RateLimiter())

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/users/refresh", handlers.RefreshToken)
	r.POST("/verify", handlers.VerifyEmail)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		// TASK 1: Get user info from JWT (Only for verified users)
		auth.GET("/users/me", middleware.VerifiedMiddleware(), handlers.GetMe)

		// TASK 2: Role-based access (admin only)
		adminOnly := auth.Group("/")
		adminOnly.Use(middleware.RoleMiddleware(models.RoleAdmin))
		{
			adminOnly.PATCH("/users/promote/:id", handlers.PromoteUser)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
