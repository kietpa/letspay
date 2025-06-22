package api

import (
	"context"
	"encoding/json"
	"fmt"
	"letspay/common/constants"
	"letspay/controller"
	"letspay/model"
	"letspay/tool/helper"
	"letspay/usecase"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type userApi struct {
	userUC   usecase.UserUsecase
	validate validator.Validate
}

func NewUserApi(
	userUC usecase.UserUsecase, validate validator.Validate,
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
		fmt.Println("register user unmarshal error: ", err)
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	if validationErrors := helper.ValidateStruct(request, a.validate); len(validationErrors) > 0 {
		fmt.Println("register user validation error: ", validationErrors)
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
		fmt.Println("login user unmarshal error: ", err)
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	if validationErrors := helper.ValidateStruct(request, a.validate); len(validationErrors) > 0 {
		fmt.Println("login user validation error: ", validationErrors)
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

	response.Status = http.StatusOK
	response.Data = userResponse

	return response, model.Error{}
}
