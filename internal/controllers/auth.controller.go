package controllers

import (
	"net/http"

	"authService/internal/models"
	"authService/internal/services"
	"authService/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AuthController struct {
	redisClient *redis.Client
}

func NewAuthController(rdb *redis.Client) *AuthController {
	return &AuthController{
		redisClient: rdb,
	}
}

// Register request payload
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest

	// Bind incoming JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service layer
	user, tokens, err := services.RegisterUser(services.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return user info and JWT tokens
	c.JSON(http.StatusOK, gin.H{
		"message":      "Signed up successfully",
		"user":         user,
		"accessToken":  tokens["accessToken"],
		"refreshToken": tokens["refreshToken"],
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest

	// Bind incoming JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service layer
	user, tokens, err := services.LoginUser(services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, ac.redisClient)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return user info and JWT tokens
	c.JSON(http.StatusOK, gin.H{
		"message":      "Loged in successfully",
		"user":         user,
		"accessToken":  tokens["accessToken"],
		"refreshToken": tokens["refreshToken"],
	})
}

func (ac *AuthController) AccessRefreshToken(c *gin.Context) {
	// refreshToken, err := c.Cookie("refreshToken")
	var body RefreshRequest
	if err := c.ShouldBindJSON(&body); err != nil || body.RefreshToken == "" {
		c.JSON(400, gin.H{"error": "refreshToken required in body"})
		return
	}
	refreshToken := body.RefreshToken

	// 1. Validate refresh token
	claims, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}
	user, err := models.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	// 3. Token version check
	if claims.TokenVersion != user.TokenVersion {
		c.JSON(401, gin.H{"error": "Refresh token expired"})
		return
	}

	err = models.IncrementTokenVersion(user.ID.Hex())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update tokenVersion"})
		return
	}
	newTokenVersion := user.TokenVersion + 1

	newAccessToken, err := utils.GenerateAccessToken(user.ID.Hex(), user.Email, user.Role, newTokenVersion)
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot create access token"})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID.Hex(), newTokenVersion)
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot create refresh token"})
		return
	}
	c.JSON(200, gin.H{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}
