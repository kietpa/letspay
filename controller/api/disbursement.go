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

type disbursementAPI struct {
	disbursementUC usecase.Disbursement
	validate       validator.Validate
}

func NewDisbursementAPI(dibursementUC usecase.Disbursement, validate validator.Validate) *disbursementAPI {
	return &disbursementAPI{
		disbursementUC: dibursementUC,
		validate:       validate,
	}
}

func (a *disbursementAPI) GetDisbursement(
	ctx context.Context, param map[string]string,
) (controller.Data, model.Error) {
	response := controller.Data{}

	disbursement, err := a.disbursementUC.GetDisbursement(ctx, param["referenceId"])
	if err.Code != 0 { // 0 = success
		return controller.Data{}, err
	}

	response.Status = http.StatusOK
	response.Data = disbursement

	return response, model.Error{}
}

func (a *disbursementAPI) CreateDisbursement(
	ctx context.Context, param map[string]string,
) (controller.Data, model.Error) {
	response := controller.Data{}
	request := model.CreateDisbursementRequest{}

	if err := json.Unmarshal([]byte(param[constants.JSON_BODY]), &request); err != nil {
		fmt.Println("create disbursement unmarshal error: ", err)
		return controller.Data{}, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		}
	}

	if validationErrors := helper.ValidateStruct(request, a.validate); len(validationErrors) > 0 {
		fmt.Println("create disbursement validation error: ", validationErrors)
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.VALIDATION_ERROR,
			Errors:  validationErrors,
		}
	}

	disbursement, err := a.disbursementUC.CreateDisbursement(ctx, request)
	if err.Code != 0 {
		return controller.Data{}, err
	}

	response.Status = http.StatusOK
	response.Data = disbursement

	return response, model.Error{}
}
