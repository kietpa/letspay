package model

import "time"

type (
	UserDetail struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	RegisterUserRequest struct {
		Name     string `json:"name" validate:"required,min=4"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	RegisterUserInput struct {
		Name           string    `json:"name"`
		Email          string    `json:"email"`
		HashedPassword string    `json:"password"`
		CreatedAt      time.Time `json:"created_at"`
	}
)
