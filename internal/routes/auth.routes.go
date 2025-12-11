package routes

import (
	"authService/internal/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authController *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.AccessRefreshToken)
		// auth.POST("/logout", authController.Logout)
		// auth.GET("/me", authController.Me)
	}
}
