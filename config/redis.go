package config

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
)

var RDB *redis.Client

func ConnectRedis() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("REDIS_HOST")
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0 // fallback
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     host,
		Username: username,
		Password: password,
		DB:       db,
	})

	// Test connection
	_, err = RDB.Ping(context.Background()).Result()
	if err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	fmt.Println("✅ Connected to Redis Cloud")
}
