package model

type CreateDisbursementRequest struct {
	UserReferenceId   string  `json:"user_reference_id" validate:"required,max=15"`
	Amount            float64 `json:"amount" validate:"required,gte=5000"`
	BankCode          string  `json:"bank_code" validate:"required"`
	BankAccountNumber string  `json:"bank_account_number" validate:"required,numeric"`
	BankAccountName   string  `json:"bank_account_name" validate:"required"`
	Description       string  `json:"description,omitempty" validate:"required"`
}
