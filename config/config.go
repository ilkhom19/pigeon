package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	username string
	password string
	salt     string
)

type Config struct {
	SMTPServer       string
	SMTPPort         string
	SMTPUsername     string
	SMTPPassword     string
	VerificationSalt string
}

func NewConfig() *Config {
	return &Config{
		SMTPServer:       "smtp.yandex.com",
		SMTPPort:         "587",
		SMTPUsername:     username,
		SMTPPassword:     password,
		VerificationSalt: salt,
	}
}

func LoadEnvs() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username = os.Getenv("SMTP_USERNAME")
	password = os.Getenv("SMTP_PASSWORD")
	salt = os.Getenv("SALT")
}
