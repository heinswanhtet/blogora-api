package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host                  string
	Port                  string
	DBUser                string
	DBPassword            string
	DBAddress             string
	DBName                string
	SECRET                string
	GO_ENV                string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	AWS_BUCKET_NAME       string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Host:                  getEnv("HOST", "localhost"),
		Port:                  getEnv("PORT", "8080"),
		DBUser:                getEnv("DB_USER", "root"),
		DBPassword:            getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:             fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                getEnv("DB_NAME", "ecom"),
		SECRET:                getEnv("SECRET", ""),
		GO_ENV:                getEnv("GO_ENV", "local"),
		AWS_ACCESS_KEY_ID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWS_SECRET_ACCESS_KEY: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWS_REGION:            getEnv("AWS_REGION", ""),
		AWS_BUCKET_NAME:       getEnv("AWS_BUCKET_NAME", ""),
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
