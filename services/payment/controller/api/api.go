package api

import (
	"letspay/services/payment/common/constants"
	"letspay/services/payment/model"
	"letspay/services/payment/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type (
	ApiModule struct {
		config          model.AppConfig
		disbursementApi *disbursementApi
	}
)

func NewAPI(
	config model.AppConfig,
	disbursementUC usecase.DisbursementUsecase,
	validate validator.Validate,
) *ApiModule {
	return &ApiModule{
		config:          config,
		disbursementApi: NewDisbursementAPI(disbursementUC, validate),
	}
}

func HandleRequests(
	cfg model.AppConfig,
	disbursementUC usecase.DisbursementUsecase,
) *mux.Router {
	validate := validator.New()
	router := mux.NewRouter().StrictSlash(true)

	apiModule := NewAPI(
		cfg,
		disbursementUC,
		*validate,
	)

	disbursement := router.PathPrefix(constants.DISBURSEMENT).Subrouter()
	disbursement.HandleFunc("/{referenceId}", apiModule.GetDisbursement)
	disbursement.HandleFunc("", apiModule.CreateDisbursement)

	callback := router.PathPrefix(constants.CALLBACK).Subrouter()
	callback.HandleFunc(constants.DISBURSEMENT+"/{provider}", apiModule.CallbackDisbursement)

	return router
}
