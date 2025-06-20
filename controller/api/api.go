package api

import (
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/database"
	"letspay/repository/provider"
	"letspay/usecase"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type (
	ApiModule struct {
		config          model.AppConfig
		disbursementAPI *disbursementAPI
	}
)

func NewAPI(
	config model.AppConfig,
	disbursementUC usecase.Disbursement,
	validate validator.Validate,
) *ApiModule {
	return &ApiModule{
		config:          config,
		disbursementAPI: NewDisbursementAPI(disbursementUC, validate),
	}
}

func HandleRequests(
	cfg model.AppConfig,
	disbursementRepo database.DisbursementRepo,
	providerRepo map[int]provider.ProviderRepo,
) {
	validate := validator.New()
	router := mux.NewRouter().StrictSlash(true)

	// router.Use(controller.LoggingMiddleware)

	disbursementUC := usecase.NewDisbursementUsecase(disbursementRepo, providerRepo)

	apiModule := NewAPI(
		cfg,
		disbursementUC,
		*validate,
	)

	router.HandleFunc(constants.DISBURSEMENT+"/{referenceId}", apiModule.GetDisbursement)
	router.HandleFunc(constants.DISBURSEMENT, apiModule.CreateDisbursement)

	log.Println("API listening on port: " + cfg.Server.Port)
	http.ListenAndServe(":"+cfg.Server.Port, router)
}
