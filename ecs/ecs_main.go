package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"quotes/config"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// Initialize Redis client
func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.ElasticacheUrl, // Адреса ElastiCache
		Password: "",                    // Redis без паролю
		DB:       0,                     // Стандартна база
	})
	log.Println("Initialized Redis client for ECS")
}

// Handler for fetching quotes
func getQuotesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Retrieve the number of quotes from the "count" parameter
	count := config.DefaultCount
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil && c > 0 {
			count = c
		}
	}

	// Fetch quote keys from Redis using SCAN
	iter := redisClient.Scan(ctx, 0, "quote:*", 0).Iterator()
	var quotes []map[string]string
	for iter.Next(ctx) {
		key := iter.Val()

		quote, err := redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to fetch quote for key %s: %v", key, err)
			continue
		}
		quotes = append(quotes, quote)

		// Stop if the required number of quotes is collected
		if len(quotes) >= count {
			break
		}
	}
	if err := iter.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to scan keys: %v", err), http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

// Health Check for ECS
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Register handlers
	http.HandleFunc("/quotes", getQuotesHandler)
	http.HandleFunc("/health", healthCheckHandler)

	// Starting server
	log.Printf("ECS server is running on %s", config.ServerAddress)
	log.Fatal(http.ListenAndServe(config.ServerAddress, nil))
}
