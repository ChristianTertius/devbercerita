package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBUrlMigration string
	SecretJwt      string
	DBHost         string
	DBUser         string
	DBName         string
	DBPassword     string
	DBPort         string
	AppEnv         string
}

func LoadConfig() (*Config, error) {
	// di production tidak ada .env, jadi ignore error
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using environtment variables!")
	}

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBUrlMigration: os.Getenv("DATABASE_URL"),
		SecretJwt:      os.Getenv("SECRET_JWT"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBUser:         getEnv("DB_USER", "root"),
		DBName:         getEnv("DB_NAME", "devbercerita"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBPassword:     getEnv("DB_PASSWORD", "jjj"),
		AppEnv:         getEnv("app_env", "development"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
