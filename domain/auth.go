package domain

import (
	"context"
	"errors"
	"fmt"
	"log"

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return data.NilAuthResponse, fmt.Errorf("hashing password error: %w", err)
	}
	user := data.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	user, err = as.userRepo.Create(ctx, user)
	if err != nil {
		return data.NilAuthResponse, fmt.Errorf("create user error: %w", err)
	}

	// TODO: gen accessToken
	accessToken := "access token"

	return data.AuthResponse{
		AccessToken: accessToken,
		User:        user,
	}, err
}

func (as *AuthService) Login(ctx context.Context, input data.LoginInput) (data.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return data.NilAuthResponse, err
	}

	var user data.User
	var err error
	if data.IsEmail(input.UsernameOrEmail) {
		user, err = as.userRepo.GetByEmail(ctx, input.UsernameOrEmail)
	} else {
		user, err = as.userRepo.GetByUsername(ctx, input.UsernameOrEmail)
	}
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			return data.NilAuthResponse, data.ErrBadCredentials
		default:
			return data.NilAuthResponse, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Println("compare hash and password error:", err)
		return data.NilAuthResponse, data.ErrBadCredentials
	}

	// TODO: gen accessToken
	accessToken := "access token"

	return data.AuthResponse{
		AccessToken: accessToken,
		User:        user,
	}, nil
}
