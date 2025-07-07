package dto

import "time"

type Bank struct {
	Id        int       `db:"id"`
	BankName  string    `db:"bank_name"`
	BankCode  string    `db:"bank_code"`
	Providers string    `db:"providers"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
