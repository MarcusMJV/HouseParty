package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file loaded")
	}
}

func GetJWTKey() string {
	return os.Getenv("SECRET_JWT_KEY")
}
