package main

import (
	"log"
	"os"

	"authService/config"
	"authService/internal/controllers"
	ratelimiter "authService/internal/middlewares"
	"authService/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment")
	}
	log.Println("MONGO_URI:", os.Getenv("MONGO_URI"))
	config.ConnectRedis()
	config.ConnectMongo()

	r := gin.Default()

	r.Use(ratelimiter.RateLimiter(config.RDB))

	authController := controllers.NewAuthController()

	routes.AuthRoutes(r, authController)

	log.Println("🚀 Server running on :8080")
	r.Run(":8080")
}
