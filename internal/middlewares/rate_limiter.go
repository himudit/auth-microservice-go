package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"authService/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := context.Background()

		// 1. Extract IP
		ip := utils.GetIP(c.Request)
		if ip == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unable to get IP"})
			c.Abort()
			return
		}

		fmt.Println("Client IP:", ip)

		key := "rate_limit:" + ip

		// 2. Check if entry exists
		exists, err := rdb.Exists(ctx, key).Result()
		if err != nil {
			c.JSON(500, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}

		if exists == 0 {
			// Create entry
			err := rdb.HSet(ctx, key, map[string]interface{}{
				"tokens":         10,
				"last_refill_ts": time.Now().Unix(),
			}).Err()

			if err != nil {
				c.JSON(500, gin.H{"error": "Redis write error"})
				c.Abort()
				return
			}
			fmt.Println("Created redis entry for:", ip)
		}

		// Continue to next handler
		c.Next()
	}
}
