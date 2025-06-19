package model

type (
	CreateDisbursementProviderOutput struct {
		ReferenceId         string  `db:"reference_id"`
		ProviderReferenceId string  `db:"provider_reference_id"`
		Status              string  `db:"status"`
		Amount              float64 `db:"amount"`
		BankCode            string  `db:"bank_code"`
		BankAccountNumber   string  `db:"bank_account_number"`
		BankAccountName     string  `db:"bank_account_name"`
		Description         string  `db:"description"`
		FailureCode         string  `db:"failure_code"`
	}
)
