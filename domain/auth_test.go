package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/startdusk/twitter/data"
	mocks "github.com/startdusk/twitter/mocks/data"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	input := data.RegisterInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	
	t.Run("can register", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{
			ID:       "user_id",
			Username: "bob",
			Email:    "bob@gmail.com",
		}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(returnUser, nil)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, nil)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(returnUser, nil)

		as := NewAuthService(&userRepo)
		resp, err := as.Register(context.Background(), input)
		require.NoError(t, err)
		require.NotEmpty(t, resp.User.ID)
		require.Equal(t, input.Email, resp.User.Email)
		require.Equal(t, input.Username, resp.User.Username)
		require.NotEmpty(t, resp.AccessToken)

		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("username taken", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(returnUser, data.ErrUsernameTaken)

		as := NewAuthService(&userRepo)
		_, err := as.Register(context.Background(), input)
		require.Error(t, err)
		require.Equal(t, err, data.ErrUsernameTaken)

		require.True(t, userRepo.AssertNotCalled(t, "GetByEmail"))
		require.True(t, userRepo.AssertNotCalled(t, "Create"))
		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("email taken", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(returnUser, nil)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, data.ErrEmailTaken)

		as := NewAuthService(&userRepo)
		_, err := as.Register(context.Background(), input)
		require.Error(t, err)
		require.Equal(t, err, data.ErrEmailTaken)

		require.True(t, userRepo.AssertNotCalled(t, "Create"))
		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("create error", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(returnUser, nil)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, nil)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(returnUser, errors.New("create user error"))

		as := NewAuthService(&userRepo)
		_, err := as.Register(context.Background(), input)
		require.Error(t, err)

		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("invalid input", func(t *testing.T) {
		var userRepo mocks.UserRepo
		as := NewAuthService(&userRepo)
		_, err := as.Register(context.Background(), data.RegisterInput{})
		require.Error(t, err)

		require.True(t, userRepo.AssertNotCalled(t, "GetByUsername"))
		require.True(t, userRepo.AssertNotCalled(t, "GetByEmail"))
		require.True(t, userRepo.AssertNotCalled(t, "Create"))
		require.True(t, userRepo.AssertExpectations(t))
	})
}
