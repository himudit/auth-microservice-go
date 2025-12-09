package main

import (
	"log"

	"authService/config"
	"authService/internal/controllers"
	ratelimiter "authService/internal/middlewares"
	"authService/internal/routes"
	"authService/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment")
	}
	config.ConnectRedis()
	config.ConnectMongo()
	services.InitCollections()
	config.LoadRSAKeys()

	r := gin.Default()

	r.Use(ratelimiter.RateLimiter(config.RDB))

	authController := controllers.NewAuthController()

	routes.AuthRoutes(r, authController)

	log.Println("🚀 Server running on :8080")
	r.Run(":8080")
}
