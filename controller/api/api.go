package api

import (
	"letspay/common/constants"
	"letspay/controller/middleware"
	"letspay/model"
	"letspay/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type (
	ApiModule struct {
		config          model.AppConfig
		disbursementApi *disbursementApi
		userApi         *userApi
		logger          zerolog.Logger
	}
)

func NewAPI(
	config model.AppConfig,
	disbursementUC usecase.DisbursementUsecase,
	userUC usecase.UserUsecase,
	validate validator.Validate,
	logger zerolog.Logger,
) *ApiModule {
	return &ApiModule{
		config:          config,
		disbursementApi: NewDisbursementAPI(disbursementUC, validate, logger),
		userApi:         NewUserApi(userUC, validate, logger),
		logger:          logger,
	}
}

func HandleRequests(
	cfg model.AppConfig,
	disbursementUC usecase.DisbursementUsecase,
	userUC usecase.UserUsecase,
	logger zerolog.Logger,
) *mux.Router {
	validate := validator.New()
	router := mux.NewRouter().StrictSlash(true)

	apiModule := NewAPI(
		cfg,
		disbursementUC,
		userUC,
		*validate,
		logger,
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
