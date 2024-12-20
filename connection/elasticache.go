package connection

import (
	"context"
	"fmt"
	"log"

	"quotes/config"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// InitCache initializes the Redis client
func InitCache() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.ElasticacheUrl, // Ensure this matches your ELASTICACHE_URL
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}

// StoreQuotes stores a slice of quotes in Redis
func StoreQuotes(quotes []Quote) error {
	ctx := context.Background()
	for i, quote := range quotes {
		key := fmt.Sprintf("quote:%d", i+1)
		err := RedisClient.HSet(ctx, key, map[string]interface{}{
			"text":   quote.Text,
			"author": quote.Author,
		}).Err()
		if err != nil {
			log.Printf("Failed to store quote %d: %v", i+1, err)
			return err
		}
	}
	return nil
}

// GetQuotes retrieves a specified number of quotes from Redis
func GetQuotes(count int) ([]Quote, error) {
	ctx := context.Background()
	keys, err := RedisClient.Keys(ctx, "quote:*").Result()
	if err != nil {
		return nil, err
	}

	quotes := []Quote{}
	for i, key := range keys {
		if i >= count {
			break
		}

		text, err := RedisClient.HGet(ctx, key, "text").Result()
		if err != nil {
			log.Printf("Failed to get text for key %s: %v", key, err)
			continue
		}

		author, err := RedisClient.HGet(ctx, key, "author").Result()
		if err != nil {
			log.Printf("Failed to get author for key %s: %v", key, err)
			author = "Unknown"
		}

		quotes = append(quotes, Quote{
			Text:   text,
			Author: author,
		})
	}

	return quotes, nil
}
