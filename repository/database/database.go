package database

import (
	"context"
	"letspay/dto"
	"letspay/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	disbursementRepo struct {
		db *pgxpool.Pool
	}

	DisbursementRepo interface {
		GetDisbursement(
			ctx context.Context, transactionId string,
		) (dto.Disbursement, error)
		CreateDisbursement(
			ctx context.Context, createDisbursementInput model.CreateDisbursementInput,
		) error
		UpdateDisbursement(
			ctx context.Context, updateDisbursementInput model.UpdateDisbursementInput,
		) error
	}
)
