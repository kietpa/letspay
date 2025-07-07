package database

import (
	"context"
	"letspay/services/payment/dto"
	"letspay/services/payment/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	disbursementRepo struct {
		db *pgxpool.Pool
	}

	bankRepo struct {
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
		GetPendingDisbursements(
			ctx context.Context,
		) ([]model.GetPendingDisbursementsOutput, error)
	}

	BankRepo interface {
		GetBankByCode(
			ctx context.Context,
			bankCode string,
		) (dto.Bank, error)
		GetAllBanks(
			ctx context.Context,
		) ([]model.BankDetail, error)
	}
)
