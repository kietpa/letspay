package database

import (
	"context"
	"letspay/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	disbursementRepo struct {
		db *pgxpool.Pool
	}

	DisbursementRepo interface {
		GetDisbursement(ctx context.Context, transactionId string) (dto.Disbursement, error)
	}
)
