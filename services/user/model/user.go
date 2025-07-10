package model

import "time"

type (
	UserDetail struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	GetUserDetail struct {
		UserId  int    `json:"user_id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Webhook string `json:"webhook"`
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

	LoginUserRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	LoginUserResponse struct {
		Token string `json:"token"`
	}

	AddWebhookRequest struct {
		UserId  int    `json:"user_id"`
		Webhook string `json:"webhook"`
	}
)
