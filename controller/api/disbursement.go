package api

import (
	"context"
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
