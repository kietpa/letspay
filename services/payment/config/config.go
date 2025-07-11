package config

import (
	"letspay/services/payment/common/constants"
	"letspay/services/payment/model"
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

	xenditProvider := model.Provider{
		ApiKey:        os.Getenv("XENDIT_API_KEY"),
		BaseUrl:       os.Getenv("XENDIT_BASE_URL"),
		CallbackToken: os.Getenv("XENDIT_CALLBACK_TOKEN"),
	}

	midtransProvider := model.Provider{
		ApiKey:  os.Getenv("MIDTRANS_SERVER_KEY"),
		BaseUrl: os.Getenv("MIDTRANS_URL"),
	}

	return model.AppConfig{
		Server: model.Server{
			Port:    "8080",
			Timeout: 30,
		},
		Provider: map[int]model.Provider{
			constants.XENDIT_PROVIDER_ID:   xenditProvider,
			constants.MIDTRANS_PROVIDER_ID: midtransProvider,
		},
		Redis: model.Redis{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		RabbitMqUrl: os.Getenv("RABBITMQ_URL"),
	}
}
