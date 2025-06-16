package model

import "time"

// request = from user
// input = from usecase to repo/DB
type (
	DisbursementDetail struct {
		ReferenceId       string
		UserReferenceId   string
		Status            string
		Amount            float64
		BankCode          string
		CreatedAt         time.Time
		BankAccountNumber string
		BankAccountName   string
		Description       string
		FailureCode       string
	}

	CreateDisbursementRequest struct {
		UserReferenceId   string  `json:"user_reference_id"`
		Amount            float64 `json:"amount"`
		BankCode          string  `json:"bank_code"`
		BankAccountNumber string  `json:"bank_account_number"`
		BankAccountName   string  `json:"bank_account_name"`
		Description       string  `json:"description,omitempty"`
	}

	CreateDisbursementInput struct {
		UserId            int
		ReferenceId       string
		UserReferenceId   string
		Status            string
		Amount            float64
		BankCode          string
		CreatedAt         time.Time
		BankAccountNumber string
		BankAccountName   string
		Description       string
	}
)
