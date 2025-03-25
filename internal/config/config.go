package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIPort         string
	EthNodeURL      string
	DBConnectionURL string
	JWTSecret       string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}

	return Config{
		APIPort:         getConfigOrFail("API_PORT"),
		EthNodeURL:      getConfigOrFail("ETH_NODE_URL"),
		DBConnectionURL: getConfigOrFail("DB_CONNECTION_URL"),
		JWTSecret:       getConfigOrFail("JWT_SECRET"),
	}
}

func getConfigOrFail(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("environment variable not found: %s", key)
	}
	return value
}
