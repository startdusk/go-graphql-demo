package data

import "errors"

var (
	ErrBadCredentials    = errors.New("username_or_email/password wrong combination")
	ErrValidation        = errors.New("validation error")
	ErrNotFound          = errors.New("not found")
	ErrGenAccessToken    = errors.New("generate access token error")
	ErrGenRefreshToken   = errors.New("generate refresh token error")
	ErrParseToken        = errors.New("cannnot parse token")
	ErrNoUserIDInContext = errors.New("no user id in context")
	ErrUnauthenticated   = errors.New("unauthenticated")
	ErrInvalidUUID       = errors.New("invalid uuid")
	ErrForbidden         = errors.New("forbidden")
)
