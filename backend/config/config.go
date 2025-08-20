package config

import (
	"os"
)

// Config 應用程式配置結構
type Config struct {
	Environment    string
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiryHours int
	RedisHost      string
	RedisPort      string
	RedisPassword  string
	FrontendURL    string
}

// LoadConfig 從環境變數載入配置
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 默認端口
	}

	jwtExpiryStr := os.Getenv("JWT_EXPIRY_HOURS")
	jwtExpiry := 24 // 默認 24 小時
	
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	return &Config{
		Environment:    getEnv("ENV", "development"),
		Port:           port,
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiryHours: jwtExpiry,
		RedisHost:      redisHost,
		RedisPort:      redisPort,
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		FrontendURL:    frontendURL,
	}
}

// getEnv 獲取環境變數，若不存在則返回默認值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
