package main

import (
	"log"

	"authService/config"
	"authService/internal/routes"
	"authService/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to Redis (global singleton)
	config.ConnectRedis()

	// Create Gin router instance
	r := gin.Default()

	// Apply global middleware (if you want)
	r.Use(ratelimiter.RateLimiter(config.RDB))

	// Register routes
	routes.AuthRoutes(r)

	log.Println("🚀 Server running on :8080")
	r.Run(":8080")
}
