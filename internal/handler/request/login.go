package request

import "github.com/go-playground/validator/v10"

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=8"`
}

func (lr *LoginRequest) Validate(valid *validator.Validate) error {
	return valid.Struct(lr)
}
