package model

type DisbursementRequestEvent struct {
	UserId            int     `json:"user_id"`
	UserReferenceId   string  `json:"user_reference_id"`
	Amount            float64 `json:"amount"`
	BankCode          string  `json:"bank_code"`
	BankAccountNumber string  `json:"bank_account_number"`
	BankAccountName   string  `json:"bank_account_name"`
	Description       string  `json:"description,omitempty"`
}

// based on model.DisbursementDetail
// should use this if need additional metadata for mq
// type DisbursementBaseEvent struct {
// 	ReferenceId       string    `json:"reference_id"`
// 	UserReferenceId   string    `json:"user_reference_id"`
// 	Status            string    `json:"status"`
// 	Amount            float64   `json:"amount"`
// 	BankCode          string    `json:"bank_code"`
// 	CreatedAt         time.Time `json:"created_at"`
// 	BankAccountNumber string    `json:"bank_account_number"`
// 	BankAccountName   string    `json:"bank_account_name"`
// 	Description       string    `json:"description"`
// 	FailureCode       string    `json:"failure_code"`
// }

type DisbursementCompletedEvent struct {
	UserId int `json:"user_id"`
	DisbursementDetail
}

type DisbursementFailedEvent struct {
	UserId int `json:"user_id"`
	DisbursementDetail
}
