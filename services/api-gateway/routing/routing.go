package routing

import (
	"letspay/pkg/auth"
	"letspay/services/api-gateway/common/constants"
	"letspay/services/api-gateway/handler"

	"github.com/gorilla/mux"
)

func InitRouting(
	router *mux.Router,
	userUrl string,
	paymentUrl string,
	apiHandler *handler.ApiHandler,
) {

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/register", handler.NewReverseProxy(userUrl))
	user.HandleFunc("/login", handler.NewReverseProxy(userUrl))

	disbursement := router.PathPrefix(constants.DISBURSEMENT).Subrouter()
	disbursement.Use(auth.AuthMiddleware)
	disbursement.HandleFunc("/{referenceId}", handler.NewReverseProxy(paymentUrl))
	disbursement.HandleFunc("", apiHandler.RequestDisbursement).Methods("POST")

	callback := router.PathPrefix(constants.CALLBACK).Subrouter()
	callback.HandleFunc(constants.DISBURSEMENT+"/{provider}", handler.NewReverseProxy(paymentUrl))

}
