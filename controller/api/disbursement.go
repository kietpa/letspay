package api

import (
	"context"
	"letspay/controller"
	"letspay/usecase"
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
) (controller.Data, error) {
	response := controller.Data{}

	return response, nil
}
