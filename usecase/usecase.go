package usecase

import (
	"context"
	"letspay/model"
	"letspay/repository/database"
)

type (
	disbursementUsecase struct {
		disbursementRepo database.DisbursementRepo
	}

	Disbursement interface {
		GetDisbursement(ctx context.Context, refid string) (model.DisbursementDetail, model.Error)
		CreateDisbursement(ctx context.Context, createDisbursementRequest model.CreateDisbursementRequest) (model.DisbursementDetail, model.Error)
	}
)
