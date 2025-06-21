package database

import (
	"context"
	"letspay/dto"
	"letspay/model"

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
	created_at
	FROM users
	WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		resp.UserId,
		resp.Name,
		resp.HashedPassword,
		resp.CreatedAt,
	)

	return resp, err
}
