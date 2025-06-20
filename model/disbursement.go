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
		UserReferenceId   string  `json:"user_reference_id" validate:"required,max=15"`
		Amount            float64 `json:"amount" validate:"required,gte=5000"`
		BankCode          string  `json:"bank_code" validate:"required"`
		BankAccountNumber string  `json:"bank_account_number" validate:"required,numeric"`
		BankAccountName   string  `json:"bank_account_name" validate:"required,alpha"`
		Description       string  `json:"description,omitempty" validate:"required"`
	}

	CreateDisbursementInput struct {
		UserId            int       `db:"user_id"`
		ReferenceId       string    `db:"reference_id"`
		UserReferenceId   string    `db:"user_reference_id"`
		Status            string    `db:"status"`
		Amount            float64   `db:"amount"`
		BankCode          string    `db:"bank_code"`
		CreatedAt         time.Time `db:"created_at"`
		BankAccountNumber string    `db:"bank_account_number"`
		BankAccountName   string    `db:"bank_account_name"`
		Description       string    `db:"description"`
	}

	UpdateDisbursementInput struct {
		ReferenceId         string    `db:"reference_id"`
		ProviderId          int       `db:"provider_id,omitempty"`
		ProviderReferenceId string    `db:"provider_reference_id,omitempty"`
		Status              string    `db:"status"`
		UpdatedAt           time.Time `db:"updated_at"`
		FailureCode         string    `db:"failure_code,omitempty"`
	}
)
