package model

type (
	CreateDisbursementProviderOutput struct {
		ReferenceId         string
		ProviderReferenceId string
		Status              string
		Amount              float64
		BankCode            string
		BankAccountNumber   string
		BankAccountName     string
		Description         string
		FailureCode         string `json:"failure_code,omitempty"`
		StatusCode          int    `json:"status_code,omitempty"`
	}
)
