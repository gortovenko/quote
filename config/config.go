package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	BaseUrl        string
	DefaultCount   int
	ServerAddress  string
	ElasticacheUrl string
	CacheProvider  string
	AwsRegion      string
	RunLocalMode   bool
)

func init() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	BaseUrl = getEnv("BASE_URL", "https://quotes.toscrape.com")
	DefaultCount = getEnvInt("DEFAULT_COUNT", 100)
	ServerAddress = getEnv("SERVER_ADDRESS", ":8080")
	ElasticacheUrl = getEnv("ELASTICACHE_URL", "redis:6379")
	CacheProvider = getEnv("CACHE_PROVIDER", "elasticache")
	AwsRegion = getEnv("AWS_REGION", "us-east-1")
	RunLocalMode = getEnvBool("RUN_LOCAL_MODE", true)
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Environment variable %s not set. Using default: %s", key, defaultVal)
		return defaultVal
	}
	return val
}

func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		log.Printf("Environment variable %s not set. Using default: %d", key, defaultVal)
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Invalid value for %s: %s. Using default: %d", key, valStr, defaultVal)
		return defaultVal
	}
	return val
}

func getEnvBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		log.Printf("Environment variable %s not set. Using default: %v", key, defaultVal)
		return defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		log.Printf("Invalid value for %s: %s. Using default: %v", key, valStr, defaultVal)
		return defaultVal
	}
	return val
}
