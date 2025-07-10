package usecase

import (
	"context"
	"letspay/pkg/db"
	"letspay/services/user/model"
	"letspay/services/user/repository/database"
)

type (
	userUsecase struct {
		userRepo  database.UserRepo
		redisRepo *db.RedisClient
	}

	UserUsecase interface {
		RegisterUser(
			ctx context.Context, registerUserRequest model.RegisterUserRequest,
		) (model.UserDetail, model.Error)
		LoginUser(
			ctx context.Context, loginUserRequest model.LoginUserRequest,
		) (model.LoginUserResponse, model.Error)
		GetUser(
			ctx context.Context, userId int,
		) (model.GetUserDetail, model.Error)
		AddWebhook(
			ctx context.Context,
			input model.AddWebhookRequest,
		) model.Error
	}
)
