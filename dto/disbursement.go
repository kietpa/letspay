package dto

import "time"

type Disbursement struct {
	Id                  int       `db:"id"`
	UserId              int       `db:"user_id"`
	ReferenceId         string    `db:"reference_id"`
	UserReferenceId     string    `db:"user_reference_id"`
	ProviderId          int       `db:"provider_id"`
	ProviderReferenceId string    `db:"provider_reference_id"`
	Status              string    `db:"status"`
	Amount              float64   `db:"amount"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
	BankAccountNumber   string    `db:"bank_account_number"`
	BankAccountName     string    `db:"bank_account_name"`
	Description         string    `db:"description"`
}
