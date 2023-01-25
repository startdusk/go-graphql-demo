package data

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
)

type RegisterInput struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

var NilAuthResponse AuthResponse

type AuthResponse struct {
	AccessToken string
	User        User
}

type AuthService interface {
	Register(ctx context.Context, user User) (AuthResponse, error)
}

func (in *RegisterInput) Sanitize() {
	in.Username = strings.TrimSpace(in.Username)
	in.Email = strings.TrimSpace(in.Email)
	in.Password = strings.TrimSpace(in.Password)
	in.ConfirmPassword = strings.TrimSpace(in.ConfirmPassword)
}

var (
	UsernameMinLength = 2
	PasswordMinLength = 6
)

func (in RegisterInput) Validate() error {
	if len(in.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username not long enough, (%d) characters as least", ErrValidation, UsernameMinLength)
	}

	if _, err := mail.ParseAddress(in.Email); err != nil {
		return fmt.Errorf("%w: email not valid", ErrValidation)
	}

	if len(in.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password not long enough, (%d) characters as least", ErrValidation, PasswordMinLength)
	}

	if in.Password != in.ConfirmPassword {
		return fmt.Errorf("%w: confirm password must match the password", ErrValidation)
	}

	return nil
}
