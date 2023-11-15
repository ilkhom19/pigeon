package config

import (
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
	username = os.Getenv("SMTP_USERNAME")
	password = os.Getenv("SMTP_PASSWORD")
	salt = os.Getenv("SALT")
}
