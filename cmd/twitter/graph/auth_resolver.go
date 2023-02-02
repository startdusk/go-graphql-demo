package graph

import (
	"context"
	"errors"

	"github.com/startdusk/twitter/data"
)

func (m *mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	resp, err := m.AuthService.Register(ctx, data.RegisterInput{
		Username:        input.Username,
		Email:           input.Email,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrValidation) ||
			errors.Is(err, data.ErrUsernameTaken) ||
			errors.Is(err, data.ErrEmailTaken):
			return nil, writeBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}
	return &AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User: &User{
			ID:        resp.User.ID,
			Username:  resp.User.Username,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
		},
	}, nil
}

func (m *mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	resp, err := m.AuthService.Login(ctx, data.LoginInput{
		UsernameOrEmail: input.UsernameOrEmail,
		Password:        input.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrValidation) ||
			errors.Is(err, data.ErrBadCredentials):
			return nil, writeBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}

	return &AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User: &User{
			ID:        resp.User.ID,
			Username:  resp.User.Username,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
		},
	}, nil
}
