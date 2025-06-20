package config

import (
	"letspay/common/constants"
	"letspay/model"
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
		log.Printf("Failed to load env variables, err=%v", err)
		panic(err)
	}

	xenditProvider := model.Provider{
		ApiKey:  os.Getenv("XENDIT_API_KEY"),
		BaseUrl: os.Getenv("XENDIT_BASE_URL"),
	}

	return model.AppConfig{
		Server: model.Server{
			Port:    "8080",
			Timeout: 30,
		},
		Provider: map[int]model.Provider{
			constants.XENDIT_PROVIDER_ID: xenditProvider,
		},
	}
}
