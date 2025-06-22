package api

import (
	"letspay/common/constants"
	"letspay/controller/middleware"
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
	disbursementRepo database.DisbursementRepo,
	userRepo database.UserRepo,
	providerRepo map[int]provider.ProviderRepo,
) {
	validate := validator.New()
	router := mux.NewRouter().StrictSlash(true)

	disbursementUC := usecase.NewDisbursementUsecase(disbursementRepo, providerRepo)
	userUC := usecase.NewUserUsecase(userRepo)

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

	log.Println("API listening on port: " + cfg.Server.Port)
	http.ListenAndServe(":"+cfg.Server.Port, router)
}
