package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"authService/internal/utils"
)

func RateLimiter(rdb *redis.Client) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := context.Background()

			// 1️⃣ Extract IP
			ip := utils.GetIP(r)
			if ip == "" {
				http.Error(w, "Unable to get IP", http.StatusForbidden)
				return
			}

			fmt.Println("Client IP:", ip) // For testing

			// 2️⃣ Create Redis entry if not exists
			key := "rate_limit:" + ip

			exists, err := rdb.Exists(ctx, key).Result()
			if err != nil {
				http.Error(w, "Redis error", 500)
				return
			}

			if exists == 0 {
				// Key does not exist → create it
				err := rdb.HSet(ctx, key, map[string]interface{}{
					"tokens":         10,                // initial tokens
					"last_refill_ts": time.Now().Unix(), // timestamp
				}).Err()
				if err != nil {
					http.Error(w, "Redis set error", 500)
					return
				}

				fmt.Println("Redis entry created for IP:", ip)
			} else {
				fmt.Println("Redis entry already exists for IP:", ip)
			}

			next.ServeHTTP(w, r)
		})
	}
}
