package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuration struct for the backend
type Config struct {
	Port          string
	DataBaseURL   string
	CloudinaryURL string
	JWTSecret     string
}

// Loads the backend configuration
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	port := os.Getenv("GOLANG_SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	dbURL := os.Getenv("NEON_DB_URI")
	if dbURL == "" {
		return nil, fmt.Errorf("NEON_DB_URI is required but not set")
	}

	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		return nil, fmt.Errorf("CLOUDINARY_URL is required but not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is required but not set")
	}

	config := &Config{
		Port:          port,
		DataBaseURL:   dbURL,
		CloudinaryURL: cloudinaryURL,
		JWTSecret:     jwtSecret,
	}

	return config, nil
}
