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

	disbursement, _ := a.disbursementUC.GetDisbursement(ctx, param["referenceId"])
	// if err.Message != "" {
	// 	log.Println("API error")

	// 	response.Status = err.Code

	// 	return response, nil //TODO: change error to custom API error
	// }

	response.Data = disbursement

	return response, nil
}
