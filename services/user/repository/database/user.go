package database

import (
	"context"
	"letspay/services/user/dto"
	"letspay/services/user/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserRepo(db *pgxpool.Pool) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) RegisterUser(
	ctx context.Context,
	registerUserInput model.RegisterUserInput,
) error {
	query := `INSERT INTO users
	(
	name,
	email,
	password,
	created_at
	)
	VALUES
	(
	$1,
	$2,
	$3,
	$4
	)`

	_, err := r.db.Exec(ctx, query,
		registerUserInput.Name,
		registerUserInput.Email,
		registerUserInput.HashedPassword,
		registerUserInput.CreatedAt,
	)

	return err
}

func (r *userRepo) GetUserByEmail(
	ctx context.Context,
	email string,
) (dto.User, error) {
	resp := dto.User{}
	query := `SELECT
	id,
	name,
	email,
	password,
	created_at,
	webhook
	FROM users
	WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&resp.UserId,
		&resp.Name,
		&resp.Email,
		&resp.HashedPassword,
		&resp.CreatedAt,
		&resp.Webhook,
	)

	return resp, err
}

func (r *userRepo) GetUserById(
	ctx context.Context,
	userId int,
) (dto.User, error) {
	resp := dto.User{}
	query := `SELECT
	id,
	name,
	email,
	password,
	created_at,
	webhook
	FROM users
	WHERE id = $1`

	err := r.db.QueryRow(ctx, query, userId).Scan(
		&resp.UserId,
		&resp.Name,
		&resp.Email,
		&resp.HashedPassword,
		&resp.CreatedAt,
		&resp.Webhook,
	)

	return resp, err
}

func (r *userRepo) UpdateUserWebhook(
	ctx context.Context,
	webhook string,
	userId int,
) error {
	query := `UPDATE users
	SET webhook = $1
	WHERE id = $2`

	_, err := r.db.Exec(ctx, query,
		webhook,
		userId,
	)

	return err
}
