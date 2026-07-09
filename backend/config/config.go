package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                 string
	AppEnv                  string
	AppDomain               string
	Port                    string
	DatabaseURL             string
	AdminHash               string
	JwtAccessSecret         string
	JwtRefreshSecret        string
	JwtAccessExpiresMinutes int
	JwtRefreshExpiresDays   int
	EncryptSecretKey        string
	Memory                  uint32
	Iterations              uint32
	Parallelism             uint8
	SaltLength              uint32
	KeyLength               uint32
}

func parseArgonParams(paramStr string) (uint32, uint32, uint8, uint32, uint32, error) {
	var memory, iterations, saltLength, keyLength uint32
	var parallelism uint8
	_, err := fmt.Sscanf(paramStr, "m=%d,t=%d,p=%d,sl=%d,kl=%d",
		&memory, &iterations, &parallelism, &saltLength, &keyLength)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	return memory, iterations, parallelism, saltLength, keyLength, nil
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

	appDomain := os.Getenv("APP_DOMAIN")
	if appDomain == "" {
		log.Fatal("❌ APP_DOMAIN is required.")
	}

	encryptSecretKey := os.Getenv("ENCRYPT_SECRET_KEY")
	if encryptSecretKey == "" {
		log.Fatal("❌ ENCRYPT_SECRET_KEY is required.")
	}

	argonParamsEnv := os.Getenv("ARGON_PARAMS")
	if argonParamsEnv == "" {
		log.Fatal("❌ ARGON_PARAMS is required in .env")
	}

	memory, iterations, parallelism, saltLength, keyLength, err := parseArgonParams(argonParamsEnv)
	if err != nil {
		log.Fatal("❌ ARGON_PARAMS format is invalid. Expected format: m=65536,t=3,p=4,sl=16,kl=32")
	}

	return &Config{
		AppName:                 appName,
		AppEnv:                  env,
		AppDomain:               appDomain,
		Port:                    port,
		DatabaseURL:             dbURL,
		AdminHash:               adminHash,
		JwtAccessSecret:         jwtAccessSecret,
		JwtRefreshSecret:        jwtRefreshSecret,
		JwtAccessExpiresMinutes: aToMinutes,
		JwtRefreshExpiresDays:   rToDays,
		EncryptSecretKey:        encryptSecretKey,
		Memory:                  memory,
		Iterations:              iterations,
		Parallelism:             parallelism,
		SaltLength:              saltLength,
		KeyLength:               keyLength,
	}
}
