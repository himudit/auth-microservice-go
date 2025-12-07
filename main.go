package main

import (
	"log"
	"net/http"
	"authService/config"
	"authService/internal/middlewares"
)

func main() {
	// Connect to Redis once
	config.ConnectRedis()

	mux := http.NewServeMux()

	// Pass global RDB to middleware
	mux.Handle("/test", middleware.RateLimiter(config.RDB)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})))

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
