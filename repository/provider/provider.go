package provider

import (
	"context"
	"letspay/model"
	"net/http"
)

type (
	ProviderRepo interface {
		ExecuteDisbursement(
			ctx context.Context, input model.CreateDisbursementInput,
		) (model.CreateDisbursementProviderOutput, error)
		GetDisbursementStatus(
			ctx context.Context, providerRefid string,
		) (model.GetDisbursementProviderResponse, error)
		ValidateCallbackToken(
			ctx context.Context, headers http.Header,
		) bool
	}
)
