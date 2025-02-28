package database

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	DisbursementRepo struct {
		db    *pgxpool.Pool
		synch sync.Mutex
	}

	Disbursement interface {
		GetDisbursementTransaction(ctx context.Context, transactionId string)
	}
)
