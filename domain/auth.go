package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/startdusk/twitter/data"
	"golang.org/x/crypto/bcrypt"
)

var enctyptPasswordCost = bcrypt.DefaultCost

type AuthService struct {
	authTokenService data.AuthTokenService
	userRepo         data.UserRepo
	refreshTokenRepo data.RefreshTokenRepo
}

func NewAuthService(ur data.UserRepo, ats data.AuthTokenService, rtr data.RefreshTokenRepo) *AuthService {
	return &AuthService{
		authTokenService: ats,
		userRepo:         ur,
		refreshTokenRepo: rtr,
	}
}

func (as *AuthService) Register(ctx context.Context, input data.RegisterInput) (data.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return data.NilAuthResponse, err
	}

	// check if username is already taken
	u1, err := as.userRepo.GetByUsername(ctx, input.Username)
	if u1 != data.NilUser || (err != nil && !errors.Is(err, data.ErrNotFound)) {
		log.Println(err)
		return data.NilAuthResponse, data.ErrUsernameTaken
	}

	// check if email is already taken
	u2, err := as.userRepo.GetByEmail(ctx, input.Email)
	if u2 != data.NilUser || (err != nil && !errors.Is(err, data.ErrNotFound)) {
		log.Println(err)
		return data.NilAuthResponse, data.ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), enctyptPasswordCost)
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

	accessToken, err := as.authTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		log.Printf("%+v", err)
		return data.NilAuthResponse, data.ErrGenAccessToken
	}

	tokenID := uuid.NewString()
	refreshToken, expiredAt, err := as.authTokenService.CreateRefreshToken(ctx, user, tokenID)
	if err != nil {
		log.Printf("%+v", err)
		return data.NilAuthResponse, data.ErrGenRefreshToken
	}
	if _, err := as.refreshTokenRepo.Create(ctx, data.CreateRefreshTokenParams{
		Sub:       user.ID,
		TokenID:   tokenID,
		ExpiredAt: expiredAt,
	}); err != nil {
		log.Printf("%+v", err)
		return data.NilAuthResponse, data.ErrCreateSession
	}

	return data.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
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

	accessToken, err := as.authTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		log.Printf("%+v", err)
		return data.NilAuthResponse, data.ErrGenAccessToken
	}

	tokenID := uuid.NewString()
	refreshToken, expiredAt, err := as.authTokenService.CreateRefreshToken(ctx, user, tokenID)
	if err != nil {
		log.Printf("%+v", err)
		return data.NilAuthResponse, data.ErrGenRefreshToken
	}
	if _, err := as.refreshTokenRepo.Create(ctx, data.CreateRefreshTokenParams{
		Sub:       user.ID,
		TokenID:   tokenID,
		ExpiredAt: expiredAt,
	}); err != nil {
		log.Printf("%+v", err)
		return data.NilAuthResponse, data.ErrCreateSession
	}

	return data.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, err
}

func (as *AuthService) RefreshToken(ctx context.Context, token string) (data.AuthResponse, error) {
	authToken, err := as.authTokenService.ParseToken(ctx, token)
	if err != nil {
		return data.NilAuthResponse, data.ErrBadCredentials
	}

	refreshToken, err := as.refreshTokenRepo.GetByTokenID(ctx, authToken.ID)
	if err != nil {
		return data.NilAuthResponse, data.ErrBadCredentials
	}
	if refreshToken.ExpiredAt.Before(time.Now()) {
		return data.NilAuthResponse, data.ErrRefreshTokenExpired
	}

	go func() {
		err := as.refreshTokenRepo.LastUsed(context.Background(), data.CreateRefreshTokenParams{
			Sub:     refreshToken.UserID,
			TokenID: refreshToken.TokenID,
		})
		if err != nil {
			log.Printf("log session[%s] last used error: %+v\n", refreshToken.TokenID, err)
		}
	}()

	user := data.User{
		ID: refreshToken.UserID,
	}
	accessToken, err := as.authTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return data.NilAuthResponse, data.ErrGenAccessToken
	}

	// TODO: gen refresh token ?

	return data.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: token,
		User:         user,
	}, err
}
