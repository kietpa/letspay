package usecase

import (
	"context"
	"letspay/model"
	"letspay/repository/database"
	"letspay/repository/provider"
	"letspay/tool/redis"
	"net/http"
)

type (
	disbursementUsecase struct {
		disbursementRepo database.DisbursementRepo
		providerRepo     map[int]provider.ProviderRepo
		bankRepo         database.BankRepo
		redisRepo        *redis.RedisClient
	}

	userUsecase struct {
		userRepo database.UserRepo
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

	UserUsecase interface {
		RegisterUser(
			ctx context.Context, registerUserRequest model.RegisterUserRequest,
		) (model.UserDetail, model.Error)
		LoginUser(
			ctx context.Context, loginUserRequest model.LoginUserRequest,
		) (model.LoginUserResponse, model.Error)
	}
)
