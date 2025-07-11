package api

import (
	"context"
	"encoding/json"
	"fmt"
	"letspay/pkg/helper"
	"letspay/pkg/logger"
	"letspay/services/user/common/constants"
	"letspay/services/user/controller"
	"letspay/services/user/model"
	"letspay/services/user/usecase"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type userApi struct {
	userUC   usecase.UserUsecase
	validate validator.Validate
}

func NewUserApi(
	userUC usecase.UserUsecase,
	validate validator.Validate,
) *userApi {
	return &userApi{
		userUC:   userUC,
		validate: validate,
	}
}

func (a *userApi) RegisterUser(
	ctx context.Context, param map[string]string,
) (controller.Data, model.Error) {
	response := controller.Data{}
	request := model.RegisterUserRequest{}

	if err := json.Unmarshal([]byte(param[constants.JSON_BODY]), &request); err != nil {
		logger.Error(ctx, fmt.Sprintf("[Register User] unmarshal error=%s", err))
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	logger.Info(ctx, fmt.Sprintf("[Register User] registering user email=%s", request.Email))

	if validationErrors := helper.ValidateStruct(request, a.validate); len(validationErrors) > 0 {
		logger.Error(ctx, fmt.Sprint("[Register User] user validation error=", validationErrors))
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.VALIDATION_ERROR,
			Errors:  validationErrors,
		}
	}

	userResponse, err := a.userUC.RegisterUser(ctx, request)
	if err.Code != 0 {
		return controller.Data{}, err
	}

	logger.Info(ctx, fmt.Sprintf("[Register User] Successfully register user=%+v", userResponse))

	response.Status = http.StatusOK
	response.Data = userResponse

	return response, model.Error{}
}

func (a *userApi) LoginUser(
	ctx context.Context, param map[string]string,
) (controller.Data, model.Error) {
	response := controller.Data{}
	request := model.LoginUserRequest{}

	if err := json.Unmarshal([]byte(param[constants.JSON_BODY]), &request); err != nil {
		logger.Error(ctx, fmt.Sprint("[Login User] unmarshal error=", err))
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	logger.Error(ctx, fmt.Sprintf("[Login User] logging in user email=%s", request.Email))

	if validationErrors := helper.ValidateStruct(request, a.validate); len(validationErrors) > 0 {
		logger.Error(ctx, fmt.Sprint("[Login User] validation error=", validationErrors))
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.VALIDATION_ERROR,
			Errors:  validationErrors,
		}
	}

	userResponse, err := a.userUC.LoginUser(ctx, request)
	if err.Code != 0 {
		return controller.Data{}, err
	}

	logger.Info(ctx, fmt.Sprintf("[Login User] Successfully login user JWT=%s", userResponse))

	response.Status = http.StatusOK
	response.Data = userResponse

	return response, model.Error{}
}

func (a *userApi) GetUser(
	ctx context.Context, param map[string]string,
) (controller.Data, model.Error) {
	response := controller.Data{}
	userId, _ := strconv.Atoi(param[constants.USER_ID])

	logger.Info(ctx, fmt.Sprintf("[Get User] Getting user userid=%d", userId))

	user, err := a.userUC.GetUser(ctx, userId)
	if err.Code != 0 { // 0 = success
		return controller.Data{}, err
	}

	logger.Info(ctx, fmt.Sprintf(
		"[Get User] Successfully retrieved user refid=%d user=%+v",
		userId,
		user,
	))

	response.Status = http.StatusOK
	response.Data = user

	return response, model.Error{}
}

func (a *userApi) AddWebhook(
	ctx context.Context, param map[string]string,
) (controller.Data, model.Error) {
	response := controller.Data{}
	request := model.AddWebhookRequest{}

	if err := json.Unmarshal([]byte(param[constants.JSON_BODY]), &request); err != nil {
		logger.Error(ctx, fmt.Sprintf("[Add Webhook] unmarshal error=%s", err))
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	userId, _ := strconv.Atoi(param[constants.USER_ID])

	logger.Info(ctx, fmt.Sprintf("[Add Webhook] adding user webhook=%s", request.Webhook))

	err := a.userUC.AddWebhook(ctx, request.Webhook, userId)
	if err.Code != 0 {
		return controller.Data{}, err
	}

	logger.Info(ctx, fmt.Sprintf("[Add Webhook] Successfully added webhook=%+v userid=%d",
		request.Webhook,
		userId,
	))

	response.Status = http.StatusOK
	response.Data = map[string]string{
		"message": "Successfully registered webhook",
	}

	return response, model.Error{}
}
