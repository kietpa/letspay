package config

import (
	"letspay/services/user/model"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func InitConfig() model.AppConfig {
	if err := loadEnv(); err != nil {
		log.Fatalf("Failed to load env variables, err=%v", err)
		panic(err)
	}

	return model.AppConfig{
		Server: model.Server{
			Port:    "8080",
			Timeout: 30,
		},
		Redis: model.Redis{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}
}
