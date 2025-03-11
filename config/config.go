package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	//"fmt"
)

// LoadEnv loads environment variables from the .env file
// i use railway to deploy
func LoadEnv() error {
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found, using Railway environment variables")
		}
	}
	return nil
}
