package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	DatabaseURL         string
	SupabaseURL         string
	SupabaseAnonKey     string
	SupabaseServiceKey  string
	GeminiAPIKey        string
	Port                string
	Environment         string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		GeminiAPIKey:       getEnv("GEMINI_API_KEY", ""),
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENV", "development"),
	}

	// Validate required fields
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required in .env file")
	}

	return config
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
