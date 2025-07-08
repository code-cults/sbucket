package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func LoadEnv(){
	if err:= godotenv.Load(); err!=nil{
		log.Println("No .env variables found using system environment variables")
	}
}

func GetEnv(key,fallback string) string{
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}