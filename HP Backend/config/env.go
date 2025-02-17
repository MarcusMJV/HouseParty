package config

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("error loading .env file")
	}
}

func GetJWTKey() string {
	return os.Getenv("SECRET_JWT_KEY")
}