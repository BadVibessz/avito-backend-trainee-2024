package _map

import "github.com/pkg/errors"

var (
	ErrNoSuchKey    = errors.New("no such key")
	ErrCannotConver = errors.New("cannot convert")
)
