package api

import (
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/database"
	"letspay/usecase"

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
) *ApiModule {
	return &ApiModule{
		config:          config,
		disbursementAPI: NewDisbursementAPI(disbursementUC),
	}
}

func HandleRequests(
	cfg model.AppConfig,
	disbursementRepo database.DisbursementRepo,
) {
	router := mux.NewRouter().StrictSlash(true)

	disbursementUC := usecase.NewDisbursementUsecase(disbursementRepo)

	apiModule := NewAPI(
		cfg,
		disbursementUC,
	)

	router.HandleFunc(constants.DISBURSEMENT+"/{referenceId}", apiModule.GetDisbursement)
}
