package main

import "errors"

var (
	ErrJwtEnvVarNotSet = errors.New("JWT_SECRET env variable not set")
)
