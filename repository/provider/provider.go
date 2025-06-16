package provider

type (
	ProviderRepo interface {
		ExecuteDisbursement()
		GetDisbursementStatus()
	}
)
