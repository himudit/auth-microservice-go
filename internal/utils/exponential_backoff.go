package utils

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

type ExponentialBackoffData struct {
	FailCount   int64 `json:"failCount"`
	NextAllowed int64 `json:"nextAllowed"`
}

func CheckBackoff(email string, rdb *redis.Client) (status string, remainingTime string, err error) {
	ctx := context.TODO()
	key := "backoff:" + email
	val, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return "allowed", "0", nil
	} else if err != nil {
		return "", "", err
	}

	var backoffData ExponentialBackoffData
	if err := json.Unmarshal([]byte(val), &backoffData); err != nil {
		return "", "", err
	}

	currentTime := time.Now().Unix()

	if currentTime < backoffData.NextAllowed {
		remaining := backoffData.NextAllowed - currentTime
		return "blocked", formatDuration(remaining), nil
	}

	return "allowed", "0", nil
}

func formatDuration(seconds int64) string {
	return time.Duration(seconds * int64(time.Second)).String()
}

func UpdateBackoff(email string, rdb *redis.Client) error {
	ctx := context.TODO()
	key := "backoff:" + email
	val, err := rdb.Get(ctx, key).Result()

	var data ExponentialBackoffData
	if err == redis.Nil || val == "" {
		data = ExponentialBackoffData{
			FailCount:   1,
			NextAllowed: time.Now().Unix() + 1, // base delay 30 seconds
		}
	} else {
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			return err
		}
		data.FailCount++
		baseDelay := 1 // seconds
		delay := float64(baseDelay) * math.Pow(2, float64(data.FailCount-1))

		// 6️⃣ Update next allowed time
		data.NextAllowed = time.Now().Unix() + int64(delay)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ttl := 15 * time.Minute
	return rdb.Set(ctx, key, jsonData, ttl).Err()
}

func ResetBackoff(email string, rdb *redis.Client) {
	rdb.Del(context.TODO(), "backoff:"+email)
}
