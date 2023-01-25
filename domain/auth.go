package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/startdusk/twitter/data"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo data.UserRepo
}

func NewAuthService(ur data.UserRepo) *AuthService {
	return &AuthService{
		userRepo: ur,
	}
}

func (as *AuthService) Register(ctx context.Context, input data.RegisterInput) (data.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return data.NilAuthResponse, err
	}

	// check if username is already taken
	if _, err := as.userRepo.GetByUsername(ctx, input.Username); err != nil && !errors.Is(err, data.ErrNotFound) {
		return data.NilAuthResponse, data.ErrUsernameTaken
	}

	// check if email is already taken
	if _, err := as.userRepo.GetByEmail(ctx, input.Email); err != nil && !errors.Is(err, data.ErrNotFound) {
		return data.NilAuthResponse, data.ErrEmailTaken
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return data.NilAuthResponse, fmt.Errorf("hashing password error: %w", err)
	}
	user := data.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashPassword),
	}

	user, err = as.userRepo.Create(ctx, user)
	if err != nil {
		return data.NilAuthResponse, fmt.Errorf("create user error: %w", err)
	}

	// gen accessToken

	return data.AuthResponse{
		AccessToken: "access token",
		User:        user,
	}, err
}
