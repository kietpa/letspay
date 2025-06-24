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
	"strconv"

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

	userId, _ := strconv.Atoi(param[constants.USER_ID])

	disbursementResponse, err := a.disbursementUC.CreateDisbursement(ctx, request, userId)
	if err.Code != 0 {
		return controller.Data{}, err
	}

	response.Status = http.StatusOK
	response.Data = disbursementResponse

	return response, model.Error{}
}

func (a *disbursementApi) CallbackDisbursement(
	ctx context.Context, param map[string]string,
) (controller.Data, model.Error) {
	response := controller.Data{}
	body := model.XenditDisbursementCallback{}
	var decodedHeaders http.Header

	if err := json.Unmarshal([]byte(param[constants.REQUEST_HEADERS]), &decodedHeaders); err != nil {
		fmt.Println("callback disbursement unmarshal error: ", err)
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	if !a.disbursementUC.CallbackValidateToken(ctx, decodedHeaders, param[constants.PROVIDER]) {
		return controller.Data{}, model.Error{
			Code:    http.StatusUnauthorized,
			Message: constants.INVALID_TOKEN_MESSAGE,
		}
	}

	// TODO: handle idempotency with redis

	if err := json.Unmarshal([]byte(param[constants.JSON_BODY]), &body); err != nil {
		fmt.Println("callback disbursement unmarshal error: ", err)
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	request := model.CallbackDisbursementRequest{
		ReferenceId: body.ExternalId,
		Status:      body.Status,
		FailureCode: body.FailureCode,
	}

	err := a.disbursementUC.CallbackDisbursement(ctx, request)
	if err.Code != 0 {
		return controller.Data{}, err
	}

	response.Status = http.StatusOK
	response.Data = "Callback received"

	return response, model.Error{}
}
