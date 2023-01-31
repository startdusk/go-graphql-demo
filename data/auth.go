package data

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

var (
	UsernameMinLength = 2
	PasswordMinLength = 6
)

var NilAuthResponse AuthResponse

var NilAuthToken AuthToken

type AuthToken struct {
	ID  string
	Sub string
}

type AuthResponse struct {
	AccessToken string
	User        User
}

type AuthTokenService interface {
	CreateRefreshToken(ctx context.Context, user User, tokenID string) (string, error)
	CreateAccessToken(ctx context.Context, user User) (string, error)
	ParseTokenFromRequest(ctx context.Context, r *http.Request) (AuthToken, error)
	ParseToken(ctx context.Context, payload string) (AuthToken, error)
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}

type RegisterInput struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

func (in *RegisterInput) Sanitize() {
	in.Username = strings.TrimSpace(in.Username)
	in.Email = strings.TrimSpace(in.Email)
	in.Password = strings.TrimSpace(in.Password)
	in.ConfirmPassword = strings.TrimSpace(in.ConfirmPassword)
}

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

type LoginInput struct {
	UsernameOrEmail string
	Password        string
}

func (in *LoginInput) Sanitize() {
	in.UsernameOrEmail = strings.TrimSpace(in.UsernameOrEmail)
	in.Password = strings.TrimSpace(in.Password)
}

func (in LoginInput) Validate() error {
	switch IsEmail(in.UsernameOrEmail) {
	case true:
		if _, err := mail.ParseAddress(in.UsernameOrEmail); err != nil {
			return fmt.Errorf("%w: email not valid", ErrValidation)
		}
	case false:
		if len(in.UsernameOrEmail) < UsernameMinLength {
			return fmt.Errorf("%w: username not long enough, (%d) characters as least", ErrValidation, UsernameMinLength)
		}
	}

	if len(in.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password not long enough, (%d) characters as least", ErrValidation, PasswordMinLength)
	}

	return nil
}

var emailRegexp = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)

func IsEmail(email string) bool {
	return emailRegexp.MatchString(email)
}
