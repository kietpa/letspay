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
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, fmt.Sprintf("[Register User] method not allowed"))
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
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)

		logger.Info(ctx, fmt.Sprintf("[Register User] request received from %s", helper.GetIP(r)))

		response, err = m.userApi.LoginUser(ctx, param)
		if err.Code != 0 {
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, fmt.Sprintf("[Login User] method not allowed"))
		controller.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	controller.RespondWithJSON(w, http.StatusOK, response)
}
