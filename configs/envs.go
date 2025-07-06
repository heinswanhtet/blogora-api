package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host       string
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
	SECRET     string
	GO_ENV     string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Host:       getEnv("HOST", "localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:     getEnv("DB_NAME", "ecom"),
		SECRET:     getEnv("SECRET", ""),
		GO_ENV:     getEnv("GO_ENV", "local"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if fallback == "" {
		panic(fmt.Sprintf("%v environment variable is missing", key))
	}

	return fallback
}
