package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Gateway struct {
	redisClient *redis.Client
	rtlScript   *redis.Script
}

type JWTpayload struct {
	Sub string `json:"sub"`
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

	// Decode the Base64URL encoded payload
	decodedBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		http.Error(w, "Failed to decode token payload", http.StatusInternalServerError)
		return
	}
	// Unmarshal the decoded JSON bytes
	var decodedPayload JWTpayload
	err = json.Unmarshal(decodedBytes, &decodedPayload)
	if err != nil {
		http.Error(w, "Failed to parse token payload", http.StatusInternalServerError)
		return
	}

	userID := decodedPayload.Sub

	redisKey := "limit:" + userID
	currentTime := time.Now().Unix()
	capacity := 20
	rps := 10

	ctx := r.Context()
	// Call Redis For Rate Limiting
	result, err := g.rtlScript.Run(ctx, g.redisClient, []string{redisKey}, capacity, rps, currentTime).Result()
	if err != nil {
		http.Error(w, "Redis execution error", http.StatusInternalServerError)
		return
	}

	// go-redis returns numeric values as int64.
	if result.(int64) == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Token Verified & Request Allowed"))
	} else {
		http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
	}
}

func main() {
	// create the redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// Ping Redis to ensure it is connected before starting the server
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis client connected and ready")
	// Load The Token Bucket Script
	scriptPath := "../redis/token_bucket.lua"
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		log.Fatalf("Failed to read script file: %v", err)
	}
	// Start the server
	app := &Gateway{
		redisClient: rdb,
		rtlScript:   redis.NewScript(string(scriptContent)),
	}
	http.HandleFunc("/", app.jwtHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
