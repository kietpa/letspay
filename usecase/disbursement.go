package usecase

import (
	"context"
	"letspay/model"
	"letspay/repository/database"
)

func NewDisbursementUsecase(
	disbursementRepo database.DisbursementRepo,
) Disbursement {
	return &disbursementUsecase{
		disbursementRepo: disbursementRepo,
	}
}

func (u disbursementUsecase) GetDisbursement(
	ctx context.Context, refid string,
) (model.DisbursementDetail, model.Error) {
	return model.DisbursementDetail{}, model.Error{}
}
