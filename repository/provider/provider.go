package provider

type (
	DisbursementProvider interface {
		ValidateAccount()
		CreateDisbursement()
		GetDisbursementStatus()
	}
)
