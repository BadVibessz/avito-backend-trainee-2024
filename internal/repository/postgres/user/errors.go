package user

import "errors"

var (
	ErrUsernameExists = errors.New("user with this username already exists")
)
