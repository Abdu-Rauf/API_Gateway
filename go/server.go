package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Gateway struct {
	redisClient *redis.Client
	rtlScript   *redis.Script
}

func (g *Gateway) jwtHandler(w http.ResponseWriter, r *http.Request) {
	// Restrict to GET method like app.get("/")
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check whether header contains Authorization token
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or Invalid Token Format", http.StatusUnauthorized)
		return
	}

	// Extract token by removing the "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Verify The Token
	payload := VerifyToken(token)
	if payload == "" {
		http.Error(w, "Invalid Token Signature", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Token Verified."))
}

func main() {
	http.HandleFunc("/", jwtHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
