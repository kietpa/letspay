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

type disbursementApi struct {
	disbursementUC usecase.DisbursementUsecase
	validate       validator.Validate
}

func NewDisbursementAPI(
	dibursementUC usecase.DisbursementUsecase, validate validator.Validate,
) *disbursementApi {
	return &disbursementApi{
		disbursementUC: dibursementUC,
		validate:       validate,
	}
}

func (a *disbursementApi) GetDisbursement(
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

func (a *disbursementApi) CreateDisbursement(
	ctx context.Context, param map[string]string,
) (controller.Data, model.Error) {
	response := controller.Data{}
	request := model.CreateDisbursementRequest{}

	if err := json.Unmarshal([]byte(param[constants.JSON_BODY]), &request); err != nil {
		fmt.Println("create disbursement unmarshal error: ", err)
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
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

	disbursementResponse, err := a.disbursementUC.CreateDisbursement(ctx, request)
	if err.Code != 0 {
		return controller.Data{}, err
	}

	response.Status = http.StatusOK
	response.Data = disbursementResponse

	return response, model.Error{}
}
