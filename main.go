package main

import (
	"log"

	"authService/config"
	"authService/internal/controllers"
	ratelimiter "authService/internal/middlewares"
	"authService/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectRedis()

	r := gin.Default()

	r.Use(ratelimiter.RateLimiter(config.RDB))

	authController := controllers.NewAuthController()

	routes.AuthRoutes(r, authController)

	log.Println("🚀 Server running on :8080")
	r.Run(":8080")
}
