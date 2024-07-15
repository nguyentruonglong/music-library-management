package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config struct to hold all configuration values
type Config struct {
	MongoHost  string
	MongoPort  string
	MongoDB    string
	Port       string
	UploadPath string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Get ENV variable to determine which .env file to load
	env := os.Getenv("ENV")
	envFile := ""

	if env == "production" {
		envFile = ".env.production"
	} else if env == "development" {
		envFile = ".env.development"
	}

	// Load .env file if exists
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("No %s file found, reading configuration from environment variables\n", envFile)
	}

	// Read environment variables
	config := &Config{
		MongoHost:  getEnv("MONGO_HOST", ""),
		MongoPort:  getEnv("MONGO_PORT", ""),
		MongoDB:    getEnv("MONGO_DB", ""),
		Port:       getEnv("PORT", ""),
		UploadPath: getEnv("UPLOAD_PATH", ""),
	}

	return config, nil
}

// getEnv gets the value of an environment variable or returns a default value if the variable is not set
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
