package banner

import "errors"

var (
	ErrNoSuchFeature = errors.New("no such feature")
	ErrNoSuchTag     = errors.New("no such tag")
)
