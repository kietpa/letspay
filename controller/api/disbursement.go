package api

import (
	"context"
	"encoding/json"
	"fmt"
	"letspay/common/constants"
	"letspay/controller"
	"letspay/model"
	"letspay/tool/helper"
	"letspay/tool/logger"
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
	dibursementUC usecase.DisbursementUsecase,
	validate validator.Validate,
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
	refid := param["referenceId"]

	logger.Info(ctx, fmt.Sprintf("[Get Disbursement] Getting disbursement refid=%s", refid))

	disbursement, err := a.disbursementUC.GetDisbursement(ctx, refid)
	if err.Code != 0 { // 0 = success
		return controller.Data{}, err
	}

	logger.Info(ctx, fmt.Sprintf(
		"[Get Disbursement] Successfully retrieved disbursement refid=%s disbursement=%+v",
		disbursement.ReferenceId,
		disbursement,
	))

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
		logger.Error(ctx, fmt.Sprint("[Create Disbursement] unmarshal error: ", err))
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	if validationErrors := helper.ValidateStruct(request, a.validate); len(validationErrors) > 0 {
		logger.Error(ctx, fmt.Sprint("[Create Disbursement] validation error: ", validationErrors))
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

	logger.Info(ctx, fmt.Sprintf("[Create Disbursement] Successfully created disbursement=%+v", disbursementResponse))

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
		logger.Error(ctx, fmt.Sprint("[Callback Disbursement] unmarshal error: ", err))
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

	if err := json.Unmarshal([]byte(param[constants.JSON_BODY]), &body); err != nil {
		logger.Error(ctx, fmt.Sprint("[Callback Disbursement] unmarshal error: ", err))
		return controller.Data{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		}
	}

	request := model.CallbackDisbursementRequest{
		ReferenceId: body.ExternalId,
		Status:      body.Status,
		FailureCode: body.FailureCode,
		WebhookId:   decodedHeaders.Get(constants.WEBHOOK_ID),
	}
	logger.Info(ctx, fmt.Sprintf(
		"[Callback Disbursement] Processing disbursement callback refid=%s body=%+v",
		request.ReferenceId,
		body,
	))

	err := a.disbursementUC.CallbackDisbursement(ctx, request)
	if err.Code != 0 {
		return controller.Data{}, err
	}

	response.Status = http.StatusOK
	response.Data = "Callback received"

	return response, model.Error{}
}
