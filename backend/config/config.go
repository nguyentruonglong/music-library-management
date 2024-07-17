package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config struct to hold all configuration values
type Config struct {
	MongoHost  string // MongoDB host address
	MongoPort  string // MongoDB port number
	MongoDB    string // MongoDB database name
	Port       string // Application server port
	UploadPath string // Path for uploaded files
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Get ENV variable to determine which .env file to load
	env := os.Getenv("ENV") // Get the ENV variable
	envFile := ""           // Initialize the envFile variable

	// Determine the .env file to load based on the ENV variable
	if env == "production" {
		envFile = ".env.production" // Load production environment variables
	} else if env == "development" {
		envFile = ".env.development" // Load development environment variables
	}

	// Load .env file if it exists
	err := godotenv.Load(envFile) // Load the specified .env file
	if err != nil {
		fmt.Printf("No %s file found, reading configuration from environment variables\n", envFile) // Print a message if the .env file is not found
	}

	// Read environment variables into the Config struct
	config := &Config{
		MongoHost:  getEnv("MONGO_HOST", ""),  // Get the value of MONGO_HOST or use the default value
		MongoPort:  getEnv("MONGO_PORT", ""),  // Get the value of MONGO_PORT or use the default value
		MongoDB:    getEnv("MONGO_DB", ""),    // Get the value of MONGO_DB or use the default value
		Port:       getEnv("PORT", ""),        // Get the value of PORT or use the default value
		UploadPath: getEnv("UPLOAD_PATH", ""), // Get the value of UPLOAD_PATH or use the default value
	}

	return config, nil // Return the loaded configuration
}

// getEnv gets the value of an environment variable or returns a default value if the variable is not set
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key) // Look up the environment variable
	if !exists {
		return defaultValue // Return the default value if the variable is not set
	}
	return value // Return the value of the environment variable
}
