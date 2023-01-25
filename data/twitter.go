package data

import "errors"

var (
	ErrBadCredentials = errors.New("username_or_email/password wrong combination")
	ErrValidation     = errors.New("validation error")
	ErrNotFound       = errors.New("not found")
)
