package api

import (
	"context"
	"io"
	"letspay/common/constants"
	"letspay/controller"
	"letspay/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var notAllowedError = model.Error{
	Code:    http.StatusMethodNotAllowed,
	Message: constants.METHOD_NOT_ALLOWED_MESSAGE,
}

func (m *ApiModule) GetDisbursement(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	// TODO: handle context
	ctx := context.Background()

	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		trxId := vars["referenceId"]
		InputParams := make(map[string]string)
		InputParams["referenceId"] = trxId
		response, err = m.disbursementApi.GetDisbursement(ctx, InputParams)
		if err.Code != 0 {
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		controller.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	controller.RespondWithJSON(w, http.StatusOK, response)

}

func (m *ApiModule) CreateDisbursement(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)
		param[constants.USER_ID] = strconv.Itoa(r.Context().Value(constants.USER_ID).(int))
		response, err = m.disbursementApi.CreateDisbursement(ctx, param)
		if err.Code != 0 {
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		controller.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	controller.RespondWithJSON(w, http.StatusOK, response)
}

func (m *ApiModule) RegisterUser(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)
		response, err = m.userApi.RegisterUser(ctx, param)
		if err.Code != 0 {
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		controller.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	controller.RespondWithJSON(w, http.StatusOK, response)
}

func (m *ApiModule) LoginUser(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)
		response, err = m.userApi.LoginUser(ctx, param)
		if err.Code != 0 {
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		controller.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	controller.RespondWithJSON(w, http.StatusOK, response)
}
