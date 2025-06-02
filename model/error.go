package model

type (
	Error struct {
		Code        string `json:"code,omitempty"`
		Message     string `json:"message,omitempty"`
		Description string `json:"description,omitempty"`
	}
)
