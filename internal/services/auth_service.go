package services

import (
	"context"
	"errors"
	"time"

	"authService/config"
	"authService/internal/models"
	"authService/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func InitCollections() {
	userCollection = config.MongoDB.Collection("users")
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // optional
}

func RegisterUser(req RegisterRequest) (*models.User, map[string]string, error) {
	ctx := context.TODO()

	// 1️⃣ Check if email already exists
	existing := userCollection.FindOne(ctx, bson.M{"email": req.Email})
	if existing.Err() == nil {
		return nil, nil, errors.New("email already exists")
	}

	// 2️⃣ Hash password using Argon2id from utils/password.go
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	// 3️⃣ Create user object
	user := &models.User{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		Email:        req.Email,
		Password:     hashedPwd,
		CreatedAt:    time.Now().Unix(),
		TokenVersion: 1, // default for new users
	}

	// 4️⃣ Insert into MongoDB
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	// 5️⃣ Generate JWT tokens (access + refresh)
	accessToken, err := utils.GenerateAccessToken(user.ID.Hex(), user.Email, user.Role, user.TokenVersion)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.Hex(), user.TokenVersion)
	if err != nil {
		return nil, nil, err
	}

	tokens := map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	return user, tokens, nil
}
