package api

import (
	"context"
	"encoding/json"
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

// @Summary		Get Disbursement
// @Description	Get disbursement transaction by reference ID
// @Tags			Disbursement
// @Accept			json
// @Produce		json
// @Param			ID				path		int		true	"reference id"
// @Param			Authorization	header		string	true	"JWT Token"
// @Success		200				{object}	[]model.DisbursementDetail
// @Failure		400				{object}	model.Error
// @Failure		401				{object}	model.Error
// @Failure		500				{object}	model.Error
// @Router			/disbursement/{referenceid} [Get]
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

// @Summary		Create Disbursement
// @Description	Create a disbursement transaction
// @Tags			Disbursement
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string							true	"JWT Token"
// @Param			data			body		model.CreateDisbursementRequest	true	"The input body"
// @Success		200				{object}	model.DisbursementDetail
// @Failure		400				{object}	model.Error
// @Failure		401				{object}	model.Error
// @Failure		500				{object}	model.Error
// @Router			/disbursement [Post]
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

// @Summary		Register
// @Description	Register a user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			data	body		model.RegisterUserRequest	true	"The input body"
// @Success		200		{object}	model.UserDetail
// @Failure		400		{object}	model.Error
// @Failure		500		{object}	model.Error
// @Router			/user/register [Post]
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

// @Summary		Login
// @Description	Login and receive an access token
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			data	body		model.LoginUserRequest	true	"The input body"
// @Success		200		{object}	model.LoginUserResponse
// @Failure		400		{object}	model.Error
// @Failure		401		{object}	model.Error
// @Failure		500		{object}	model.Error
// @Router			/user/login [Post]
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

func (m *ApiModule) CallbackDisbursement(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		vars := mux.Vars(r)
		param[constants.PROVIDER] = vars[constants.PROVIDER]

		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)
		encodedHeaders, _ := json.Marshal(r.Header)
		param[constants.REQUEST_HEADERS] = string(encodedHeaders)

		response, err = m.disbursementApi.CallbackDisbursement(ctx, param)
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
