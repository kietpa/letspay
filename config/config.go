package config

import (
	"letspay/common/constants"
	"letspay/model"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load(".sandbox.env")
	if err != nil {
		return err
	}
	return nil
}

func InitConfig() model.AppConfig {
	// env := strings.ToLower(os.Getenv("APP_ENV"))
	if err := loadEnv(); err != nil {
		log.Printf("Failed to load env variables, err=%v", err)
		panic(err)
	}

	brickProvider := model.Provider{
		ClientId:     os.Getenv("BRICK_CLIENT_ID"),
		ClientSecret: os.Getenv("BRICK_CLIENT_SECRET"),
		BaseUrl:      os.Getenv("BRICK_BASE_URL"),
	}

	return model.AppConfig{
		Provider: map[int]model.Provider{
			constants.BRICK_PROVIDER_ID: brickProvider,
		},
	}
}
