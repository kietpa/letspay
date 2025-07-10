package model

import "time"

// request = from user
// input = from usecase to repo/json
type (
	DisbursementDetail struct {
		ReferenceId       string    `json:"reference_id"`
		UserReferenceId   string    `json:"user_reference_id"`
		Status            string    `json:"status"`
		Amount            float64   `json:"amount"`
		BankCode          string    `json:"bank_code"`
		CreatedAt         time.Time `json:"created_at"`
		BankAccountNumber string    `json:"bank_account_number"`
		BankAccountName   string    `json:"bank_account_name"`
		Description       string    `json:"description"`
		FailureCode       string    `json:"failure_code"`
	}

	GetDisbursementProviderResponse struct {
		ProviderReferenceId string
		Status              string
		FailureReason       string
	}

	CreateDisbursementRequest struct {
		UserReferenceId   string  `json:"user_reference_id" validate:"required,max=15"`
		Amount            float64 `json:"amount" validate:"required,gte=5000"`
		BankCode          string  `json:"bank_code" validate:"required"`
		BankAccountNumber string  `json:"bank_account_number" validate:"required,numeric"`
		BankAccountName   string  `json:"bank_account_name" validate:"required"`
		Description       string  `json:"description,omitempty" validate:"required"`
	}

	CreateDisbursementInput struct {
		UserId            int       `json:"user_id"`
		ReferenceId       string    `json:"reference_id"`
		UserReferenceId   string    `json:"user_reference_id"`
		Status            string    `json:"status"`
		Amount            float64   `json:"amount"`
		BankCode          string    `json:"bank_code"`
		CreatedAt         time.Time `json:"created_at"`
		BankAccountNumber string    `json:"bank_account_number"`
		BankAccountName   string    `json:"bank_account_name"`
		Description       string    `json:"description"`
	}

	UpdateDisbursementInput struct {
		ReferenceId         string    `json:"reference_id"`
		ProviderId          int       `json:"provider_id,omitempty"`
		ProviderReferenceId string    `json:"provider_reference_id,omitempty"`
		Status              string    `json:"status"`
		UpdatedAt           time.Time `json:"updated_at"`
		FailureCode         string    `json:"failure_code,omitempty"`
	}

	CallbackDisbursementRequest struct {
		ReferenceId string `json:"reference_id"`
		Status      string `json:"status"`
		FailureCode string `json:"failure_code,omitempty"`
		WebhookId   string `json:"webhook-id"`
	}

	XenditDisbursementCallback struct {
		Id                      string `json:"id"`
		Created                 string `json:"created"`
		Updated                 string `json:"updated"`
		ExternalId              string `json:"external_id"`
		UserId                  string `json:"user_id"`
		BankCode                string `json:"bank_code"`
		AccountHolderName       string `json:"account_holder_name"`
		Amount                  uint32 `json:"amount"`
		DisbursementDescription string `json:"disbursement_description"`
		Status                  string `json:"status"`
		FailureCode             string `json:"failure_code,omitempty"`
		IsInstant               bool   `json:"is_instant"`
	}

	GetPendingDisbursementsOutput struct {
		ReferenceId         string `json:"reference_id" db:"reference_id"`
		ProviderReferenceId string `json:"provider_reference_id" db:"provider_reference_id"`
		ProviderId          int    `json:"provider_id" db:"provider_id"`
	}
)
