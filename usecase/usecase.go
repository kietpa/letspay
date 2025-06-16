package usecase

import (
	"context"
	"letspay/model"
	"letspay/repository/database"
	"letspay/repository/provider"
)

type (
	disbursementUsecase struct {
		disbursementRepo database.DisbursementRepo
		providerRepo     map[int]provider.ProviderRepo
	}

	Disbursement interface {
		GetDisbursement(ctx context.Context, refid string) (model.DisbursementDetail, model.Error)
		CreateDisbursement(ctx context.Context, createDisbursementRequest model.CreateDisbursementRequest) (model.DisbursementDetail, model.Error)
	}
)
