package api

import (
	"context"
	"encoding/json"
	"letspay/common/constants"
	"letspay/controller"
	"letspay/model"
	"letspay/usecase"
	"net/http"
)

type disbursementAPI struct {
	disbursementUC usecase.Disbursement
}

func NewDisbursementAPI(dibursementUC usecase.Disbursement) *disbursementAPI {
	return &disbursementAPI{
		disbursementUC: dibursementUC,
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
		// TODO: log fail to unmarshal
		return controller.Data{}, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
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
