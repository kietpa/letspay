package model

import "time"

type (
	DisbursementDetail struct {
		ReferenceId       string
		UserReferenceId   string
		Status            string
		Amount            float64
		CreatedAt         time.Time
		BankAccountNumber string
		BankAccountName   string
		Description       string
	}
)
