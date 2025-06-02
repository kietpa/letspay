package database

import (
	"context"
	"letspay/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDisbursementRepo(db *pgxpool.Pool) DisbursementRepo {
	return &disbursementRepo{
		db: db,
	}
}

func (r *disbursementRepo) GetDisbursement(ctx context.Context, transactionId string) (dto.Disbursement, error) {
	resp := dto.Disbursement{}
	query := `SELECT 
	id,
	user_id,
	reference_id,
	user_reference_id,
	provider,
	provider_response_id,
	status,
	amount,
	created_at,
	updated_at,
	bank_account_number,
	bank_account_name,
	description
	FROM disbursements
	where reference_id = $1`

	err := r.db.QueryRow(context.Background(), query, transactionId).Scan(&resp)

	return resp, err
}
