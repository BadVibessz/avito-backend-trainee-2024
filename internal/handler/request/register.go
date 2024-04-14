package request

import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=1"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func (rr *RegisterRequest) Validate(valid *validator.Validate) error { return valid.Struct(rr) }
