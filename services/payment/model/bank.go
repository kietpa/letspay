package model

type (
	BankDetail struct {
		BankName  string `json:"bank_name"`
		BankCode  string `json:"bank_code"`
		Providers string `json:"providers"`
	}

	BankCodes struct {
		BankName string `json:"bank_name"`
		BankCode string `json:"bank_code"`
	}
)
