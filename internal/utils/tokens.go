package utils

import (
	"fmt"
	"time"

	"authService/config"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenExpiry  = time.Minute * 15   // 15 mins
	RefreshTokenExpiry = time.Hour * 24 * 7 // 7 days
)

type JWTClaims struct {
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TokenVersion int    `json:"tokenVersion"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId, email, role string, tokenVersion int) (string, error) {
	claims := &JWTClaims{
		UserID:       userId,
		Email:        email,
		Role:         role,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(config.PrivateKey)
}

func GenerateRefreshToken(userId string, tokenVersion int) (string, error) {
	claims := JWTClaims{
		UserID:       userId,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(config.PrivateKey)
}

func VerifyAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return config.PublicKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func VerifyRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return config.PublicKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
