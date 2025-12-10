// Package config handles application configuration from environment variables and defaults.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the application configuration values
type Config struct {
	DBPath      string
	APIKey      string
	PlayerTag   string
	APIBaseURL  string
}

// Load loads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load() // Ignore errors if .env file doesn't exist

	cfg := &Config{
		// Database configuration
		DBPath: getEnvOrDefault("DB_PATH", "battles.db"),
		
		// API configuration
		APIBaseURL: getEnvOrDefault("API_BASE_URL", "https://api.clashroyale.com"),
		
		// Required environment variables
		APIKey: getEnv("APIKEY"),
		PlayerTag: getPlayerTag(),
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns an empty string if not set
func getEnv(key string) string {
	return os.Getenv(key)
}

// getEnvOrDefault retrieves an environment variable or returns a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getPlayerTag retrieves the player tag with validation
func getPlayerTag() string {
	playerTag := os.Getenv("PLAYERTAG")
	if playerTag == "" {
		return ""
	}

	// Ensure player tag has proper format for API request
	if playerTag[0] != '#' {
		playerTag = "#" + playerTag
	}

	return playerTag
}