package api

import (
	"letspay/services/user/common/constants"
	"letspay/services/user/model"
	"letspay/services/user/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type (
	ApiModule struct {
		config  model.AppConfig
		userApi *userApi
	}
)

func NewAPI(
	config model.AppConfig,
	userUC usecase.UserUsecase,
	validate validator.Validate,
) *ApiModule {
	return &ApiModule{
		config:  config,
		userApi: NewUserApi(userUC, validate),
	}
}

func HandleRequests(
	cfg model.AppConfig,
	userUC usecase.UserUsecase,
) *mux.Router {
	validate := validator.New()
	router := mux.NewRouter().StrictSlash(true)

	apiModule := NewAPI(
		cfg,
		userUC,
		*validate,
	)

	user := router.PathPrefix(constants.USER).Subrouter()
	user.HandleFunc("/register", apiModule.RegisterUser)
	user.HandleFunc("/login", apiModule.LoginUser)

	user.HandleFunc("/{userId}", apiModule.GetUser)
	user.HandleFunc("/webhook", apiModule.AddWebhook)

	return router
}
