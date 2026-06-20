package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	AdminHash   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env not found. Using environmental variables.")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("❌ APP_ENV is required.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if string(port[0]) != ":" {
		port = ":" + port
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL is required.")
	}

	adminHash := os.Getenv("ADMIN_SECRET_HASH")
	if adminHash == "" {
		log.Fatal("❌ ADMIN_SECRET_HASH is required.")
	}

	return &Config{
		AppEnv:      env,
		Port:        port,
		DatabaseURL: dbURL,
		AdminHash:   adminHash,
	}
}
