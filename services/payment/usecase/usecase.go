package usecase

import (
	"context"
	"letspay/pkg/db"
	"letspay/services/payment/model"
	"letspay/services/payment/repository/database"
	"letspay/services/payment/repository/provider"
	"net/http"
)

type (
	disbursementUsecase struct {
		disbursementRepo database.DisbursementRepo
		providerRepo     map[int]provider.ProviderRepo
		bankRepo         database.BankRepo
		redisRepo        *db.RedisClient
	}

	DisbursementUsecase interface {
		GetDisbursement(
			ctx context.Context, refid string,
		) (model.DisbursementDetail, model.Error)
		CreateDisbursement(
			ctx context.Context, createDisbursementRequest model.CreateDisbursementRequest, userId int,
		) (model.DisbursementDetail, model.Error)
		CallbackDisbursement(
			ctx context.Context, callbackDisbursementRequest model.CallbackDisbursementRequest,
		) model.Error
		CallbackValidateToken(
			ctx context.Context, headers http.Header, provider string,
		) bool
		CheckAndUpdatePendingDisbursements(
			ctx context.Context,
		) (int, error)
	}
)
