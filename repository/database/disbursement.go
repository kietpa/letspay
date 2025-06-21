package database

import (
	"context"
	"letspay/dto"
	"letspay/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDisbursementRepo(db *pgxpool.Pool) DisbursementRepo {
	return &disbursementRepo{
		db: db,
	}
}

func (r *disbursementRepo) GetDisbursement(
	ctx context.Context, transactionId string,
) (dto.Disbursement, error) {
	resp := dto.Disbursement{}
	query := `SELECT 
	id,
	user_id,
	reference_id,
	user_reference_id,
	provider_id,
	provider_reference_id,
	status,
	amount,
	bank_code,
	created_at,
	updated_at,
	bank_account_number,
	bank_account_name,
	description,
	failure_code
	FROM disbursements
	WHERE reference_id = $1`

	err := r.db.QueryRow(ctx, query, transactionId).Scan(
		&resp.Id,
		&resp.UserId,
		&resp.ReferenceId,
		&resp.UserReferenceId,
		&resp.ProviderId,
		&resp.ProviderReferenceId,
		&resp.Status,
		&resp.Amount,
		&resp.BankCode,
		&resp.CreatedAt,
		&resp.UpdatedAt,
		&resp.BankAccountNumber,
		&resp.BankAccountName,
		&resp.Description,
		&resp.FailureCode,
	)

	return resp, err
}

func (r *disbursementRepo) CreateDisbursement(
	ctx context.Context, createDisbursementInput model.CreateDisbursementInput,
) error {
	// this is the initial disbursement record
	// failure_code, updated_at, provider_id, provider_reference_id is null
	query := ` INSERT INTO disbursements
	(
	user_id, 
	reference_id, 
	user_reference_id, 
	status,
	amount,
	bank_code,
	created_at,
	bank_account_number,
	bank_account_name,
	description
	)
	VALUES 
	(
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10
	)`

	_, err := r.db.Exec(ctx, query,
		createDisbursementInput.UserId,
		createDisbursementInput.ReferenceId,
		createDisbursementInput.UserReferenceId,
		createDisbursementInput.Status,
		createDisbursementInput.Amount,
		createDisbursementInput.BankCode,
		createDisbursementInput.CreatedAt,
		createDisbursementInput.BankAccountNumber,
		createDisbursementInput.BankAccountName,
		createDisbursementInput.Description,
	)

	return err
}

func (r *disbursementRepo) UpdateDisbursement(
	ctx context.Context, updateDisbursementInput model.UpdateDisbursementInput,
) error {
	query := `UPDATE disbursements
	SET provider_id = $1,
	provider_reference_id = $2,
	status = $3,
	updated_at = $4,
	failure_code = $5
	WHERE reference_id = $6`

	_, err := r.db.Exec(ctx, query,
		updateDisbursementInput.ProviderId,
		updateDisbursementInput.ProviderReferenceId,
		updateDisbursementInput.Status,
		updateDisbursementInput.UpdatedAt,
		updateDisbursementInput.FailureCode,
		updateDisbursementInput.ReferenceId,
	)

	return err
}
