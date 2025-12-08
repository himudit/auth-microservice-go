package ratelimiter

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"authService/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiterData struct {
	Tokens       float64 `json:"tokens"`
	LastRefillTs int64   `json:"last_refill_ts"`
}

func RateLimiter(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		ip := utils.GetIP(c.Request)
		if ip == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to identify IP"})
			c.Abort()
			return
		}
		key := "rate_limit:" + ip
		ctx := context.Background()

		val, err := rdb.HGetAll(ctx, key).Result()

		var data RateLimiterData

		if err != nil || len(val) == 0 {
			data = RateLimiterData{
				Tokens:       10, // max tokens
				LastRefillTs: time.Now().Unix(),
			}
		} else {
			data.Tokens, _ = strconv.ParseFloat(val["tokens"], 64)
			data.LastRefillTs, _ = strconv.ParseInt(val["last_refill_ts"], 10, 64)
		}

		currentTime := time.Now().Unix()
		newTokens := float64(currentTime-data.LastRefillTs) / 6.0
		if newTokens > 0 {
			data.Tokens += newTokens
			if data.Tokens > 10 { // cap tokens
				data.Tokens = 10
			}
		}

		if data.Tokens < 1 {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		data.Tokens -= 1
		data.LastRefillTs = currentTime

		rdb.HSet(ctx, key, map[string]interface{}{
			"tokens":         data.Tokens,
			"last_refill_ts": data.LastRefillTs,
		})

		c.Next()

	}
}
