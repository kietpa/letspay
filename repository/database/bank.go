package database

import (
	"context"
	"database/sql"
	"letspay/dto"
	"letspay/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewBankRepo(db *pgxpool.Pool) BankRepo {
	return &bankRepo{
		db: db,
	}
}

func (r *bankRepo) GetBankByCode(
	ctx context.Context,
	bankCode string,
) (dto.Bank, error) {
	resp := dto.Bank{}
	var nullTime sql.NullTime

	query := `SELECT
	id,
	bank_name,
	bank_code,
	providers,
	created_at,
	updated_at
	FROM banks
	WHERE bank_code = $1`

	err := r.db.QueryRow(ctx, query, bankCode).Scan(
		&resp.Id,
		&resp.BankName,
		&resp.BankCode,
		&resp.Providers,
		&resp.CreatedAt,
		&nullTime,
	)

	if nullTime.Valid {
		resp.UpdatedAt = nullTime.Time
	}

	return resp, err
}

func (r *bankRepo) GetAllBanks(
	ctx context.Context,
) ([]model.BankDetail, error) {
	query := `SELECT
	bank_name,
	bank_code
	FROM banks`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []model.BankDetail{}, err
	}
	defer rows.Close()

	var res []model.BankDetail
	for rows.Next() {
		var out model.BankDetail
		err := rows.Scan(&out.BankName, &out.BankCode)
		if err != nil {
			return []model.BankDetail{}, err
		}
		res = append(res, out)
	}

	if err := rows.Err(); err != nil {
		return []model.BankDetail{}, err
	}

	return res, nil
}
