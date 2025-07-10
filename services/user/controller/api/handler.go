package api

import (
	"context"
	"fmt"
	"io"
	"letspay/pkg/helper"
	"letspay/pkg/logger"
	"letspay/pkg/util"
	"letspay/services/user/common/constants"
	"letspay/services/user/controller"
	"letspay/services/user/model"
	"net/http"

	"github.com/gorilla/mux"
)

var notAllowedError = model.Error{
	Code:    http.StatusMethodNotAllowed,
	Message: constants.METHOD_NOT_ALLOWED_MESSAGE,
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
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)

		logger.Info(ctx, fmt.Sprintf("[Register User] request received from %s", helper.GetIP(r)))

		response, err = m.userApi.RegisterUser(ctx, param)
		if err.Code != 0 {
			helper.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, "[Register User] method not allowed")
		helper.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
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
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)

		response, err = m.userApi.LoginUser(ctx, param)
		if err.Code != 0 {
			helper.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, "[Register User] method not allowed")
		helper.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}

func (m *ApiModule) GetUser(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	switch r.Method {
	case http.MethodGet:
		param := make(map[string]string)
		vars := mux.Vars(r)
		param[constants.USER_ID] = vars["userId"]

		response, err = m.userApi.GetUser(ctx, param)
		if err.Code != 0 {
			helper.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, "[Get User] method not allowed")
		helper.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}

func (m *ApiModule) AddWebhook(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)

		response, err = m.userApi.AddWebhook(ctx, param)
		if err.Code != 0 {
			helper.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, "[Get User] method not allowed")
		helper.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}
