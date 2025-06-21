package main

import (
	"letspay/common/constants"
	"letspay/config"
	"letspay/controller/api"
	"letspay/repository/database"
	"letspay/repository/provider"
	"letspay/repository/provider/xendit"
)

func main() {
	// config
	cfg := config.InitConfig()

	// logger
	// redis & DB
	db := config.InitDB()

	// init repo with DB instance
	disbursementRepo := database.NewDisbursementRepo(db)
	userRepo := database.NewUserRepo(db)

	// TODO: init providers (agents in test)
	// provider mapper
	xenditRepo := xendit.NewProviderRepo(
		xendit.NewProviderRepoInput{
			BaseUrl: cfg.Provider[constants.XENDIT_PROVIDER_ID].BaseUrl,
			ApiKey:  cfg.Provider[constants.XENDIT_PROVIDER_ID].ApiKey,
		},
	)

	providerRepo := map[int]provider.ProviderRepo{
		constants.XENDIT_PROVIDER_ID: xenditRepo,
	}

	// scheduler
	// mssg queue

	// routing/handler
	api.HandleRequests(cfg, disbursementRepo, userRepo, providerRepo)

	db.Close()
}
