package dto

import "time"

type User struct {
	UserId         int       `db:"id"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"password"`
	CreatedAt      time.Time `db:"created_at"`
}
