package config

import (
	"log"
	"github.com/joho/godotenv"
)

func LoadEnvironment() {
	if err := godotenv.Load(); err != nil {
		log.Println("[ENV] .env not found")
	}
}