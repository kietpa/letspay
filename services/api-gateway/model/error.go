package model

import (
	"letspay/pkg/helper"
)

type (
	Error struct {
		Code    int                      `json:"code,omitempty"`
		Message string                   `json:"message,omitempty"`
		Errors  []helper.ValidationError `json:"errors,omitempty"`
	}
)
