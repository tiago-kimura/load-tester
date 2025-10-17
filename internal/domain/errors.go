package domain

import "errors"

var (
	ErrInvalidConfig  = errors.New("invalid test configuration")
	ErrRequestTimeout = errors.New("request timeout")
	ErrRequestFailed  = errors.New("request failed")
)
