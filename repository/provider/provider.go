package provider

import (
	"context"
	"letspay/model"
)

type (
	ProviderRepo interface {
		ExecuteDisbursement(
			ctx context.Context, input model.CreateDisbursementInput,
		) (model.CreateDisbursementProviderOutput, error)
		GetDisbursementStatus(
			ctx context.Context,
		) error
	}
)
