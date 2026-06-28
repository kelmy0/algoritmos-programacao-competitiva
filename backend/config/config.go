package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                 string
	AppEnv                  string
	Port                    string
	DatabaseURL             string
	AdminHash               string
	JwtAccessSecret         string
	JwtRefreshSecret        string
	JwtAccessExpiresMinutes int
	JwtRefreshExpiresDays   int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env not found. Using environmental variables.")
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		log.Fatal("❌ APP_NAME is required.")
	}

	env := os.Getenv("APP_ENV")
	if env != "production" && env != "development" {
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

	jwtAccessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if jwtAccessSecret == "" {
		log.Fatal("❌ JWT_ACCESS_SECRET is required.")
	}

	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtRefreshSecret == "" {
		log.Fatal("❌ JWT_REFRESH_SECRET is required.")
	}

	jwtAccessExpiresMinutes := os.Getenv("JWT_ACCESS_EXPIRES_MINUTES")
	aToMinutes, err := strconv.Atoi(jwtAccessExpiresMinutes)
	if jwtAccessExpiresMinutes == "" || err != nil || aToMinutes == 0 {
		log.Fatal("❌ JWT_ACCESS_EXPIRES_MINUTES is required.")
	}

	jwtRefreshExpiresDays := os.Getenv("JWT_REFRESH_EXPIRES_DAYS")
	rToDays, err := strconv.Atoi(jwtRefreshExpiresDays)
	if jwtRefreshExpiresDays == "" || err != nil || rToDays == 0 {
		log.Fatal("❌ JWT_REFRESH_EXPIRES_DAYS is required.")
	}

	return &Config{
		AppName:                 appName,
		AppEnv:                  env,
		Port:                    port,
		DatabaseURL:             dbURL,
		AdminHash:               adminHash,
		JwtAccessSecret:         jwtAccessSecret,
		JwtRefreshSecret:        jwtRefreshSecret,
		JwtAccessExpiresMinutes: aToMinutes,
		JwtRefreshExpiresDays:   rToDays,
	}
}
