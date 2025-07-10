package database

import (
	"context"
	"letspay/services/user/dto"
	"letspay/services/user/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	userRepo struct {
		db *pgxpool.Pool
	}

	UserRepo interface {
		RegisterUser(
			ctx context.Context,
			registerUserInput model.RegisterUserInput,
		) error
		GetUserByEmail(
			ctx context.Context,
			email string,
		) (dto.User, error)
		GetUserById(
			ctx context.Context,
			userId int,
		) (dto.User, error)
		UpdateUserWebhook(
			ctx context.Context,
			input model.AddWebhookRequest,
		) error
	}
)
