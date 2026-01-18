package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerName       string
	PlazaNetName     string
	JWTSecret        string
	Port             string
	AuthSecure       bool
	HTTPOnly         bool
	PlazaNetEnabled  bool
	PlazaNetURL      string
	GamePlazaEnabled bool
}

var Cfg *Config

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found - using defaults.")
	}

	Cfg = &Config{
		ServerName:       GetEnvOrDefault("SERVER_NAME", "PlazaNet Accounts"),
		PlazaNetName:     GetEnvOrDefault("PLAZANET_NAME", "PlazaNet"),
		JWTSecret:        GetEnvOrDefault("JWT_SECRET", "anotsosecretjwtsecret"),
		Port:             GetEnvOrDefault("PORT", ":7592"),
		AuthSecure:       os.Getenv("AUTH_SECURE") == "true",
		HTTPOnly:         os.Getenv("HTTP_ONLY") == "true",
		PlazaNetEnabled:  os.Getenv("PLAZANET_ENABLED") == "true",
		PlazaNetURL:      GetEnvOrDefault("PLAZANET_URL", "https://app.plazanet.org"),
		GamePlazaEnabled: os.Getenv("GAMEPLAZA_ENABLED") == "true",
	}

	return Cfg
}
