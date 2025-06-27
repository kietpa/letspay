package api

import (
	"letspay/common/constants"
	"letspay/controller/middleware"
	"letspay/model"
	"letspay/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type (
	ApiModule struct {
		config          model.AppConfig
		disbursementApi *disbursementApi
		userApi         *userApi
	}
)

func NewAPI(
	config model.AppConfig,
	disbursementUC usecase.DisbursementUsecase,
	userUC usecase.UserUsecase,
	validate validator.Validate,
) *ApiModule {
	return &ApiModule{
		config:          config,
		disbursementApi: NewDisbursementAPI(disbursementUC, validate),
		userApi:         NewUserApi(userUC, validate),
	}
}

func HandleRequests(
	cfg model.AppConfig,
	disbursementUC usecase.DisbursementUsecase,
	userUC usecase.UserUsecase,
) *mux.Router {
	validate := validator.New()
	router := mux.NewRouter().StrictSlash(true)

	apiModule := NewAPI(
		cfg,
		disbursementUC,
		userUC,
		*validate,
	)

	user := router.PathPrefix(constants.USER).Subrouter()
	user.HandleFunc("/register", apiModule.RegisterUser)
	user.HandleFunc("/login", apiModule.LoginUser)

	disbursement := router.PathPrefix(constants.DISBURSEMENT).Subrouter()
	disbursement.Use(middleware.AuthMiddleware)
	disbursement.HandleFunc("/{referenceId}", apiModule.GetDisbursement)
	disbursement.HandleFunc("", apiModule.CreateDisbursement)

	callback := router.PathPrefix(constants.CALLBACK).Subrouter()
	callback.HandleFunc(constants.DISBURSEMENT+"/{provider}", apiModule.CallbackDisbursement)

	return router
}
