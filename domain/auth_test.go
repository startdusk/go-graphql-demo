package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/startdusk/twitter/data"
	"github.com/startdusk/twitter/faker"
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
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(data.NilUser, nil)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(data.NilUser, nil)
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

func TestAuthService_Login(t *testing.T) {
	password := "password"
	usernameInput := data.LoginInput{
		UsernameOrEmail: "bob",
		Password:        password,
	}
	emailInput := data.LoginInput{
		UsernameOrEmail: "bob@gmail.com",
		Password:        password,
	}

	t.Run("username can login", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{
			ID:       "user_id",
			Username: "bob",
			Email:    "bob@gmail.com",
			Password: faker.HashedPassword,
		}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(returnUser, nil)

		as := NewAuthService(&userRepo)
		resp, err := as.Login(context.Background(), usernameInput)
		require.NoError(t, err)
		require.NotEmpty(t, resp.User.ID)
		require.Equal(t, usernameInput.UsernameOrEmail, resp.User.Username)
		require.NotEmpty(t, resp.AccessToken)

		require.True(t, userRepo.AssertNotCalled(t, "GetByEmail"))
		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("email can login", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{
			ID:       "user_id",
			Username: "bob",
			Email:    "bob@gmail.com",
			Password: faker.HashedPassword,
		}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, nil)

		as := NewAuthService(&userRepo)
		resp, err := as.Login(context.Background(), emailInput)
		require.NoError(t, err)
		require.NotEmpty(t, resp.User.ID)
		require.Equal(t, emailInput.UsernameOrEmail, resp.User.Email)
		require.NotEmpty(t, resp.AccessToken)

		require.True(t, userRepo.AssertNotCalled(t, "GetByUsername"))
		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("user not found", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, data.ErrNotFound)
		as := NewAuthService(&userRepo)
		_, err := as.Login(context.Background(), emailInput)
		require.Error(t, err)
		require.Equal(t, err, data.ErrBadCredentials)

		require.True(t, userRepo.AssertNotCalled(t, "GetByUsername"))
		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("internal error", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		expectedErr := errors.New("internal error")
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, expectedErr)
		as := NewAuthService(&userRepo)
		_, err := as.Login(context.Background(), emailInput)
		require.Error(t, err)
		require.Equal(t, err, expectedErr)

		require.True(t, userRepo.AssertNotCalled(t, "GetByUsername"))
		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("user password error", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{
			ID:       "user_id",
			Username: "bob",
			Email:    "bob@gmail.com",
			Password: faker.HashedPassword,
		}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, nil)

		emailInput.Password = "invalid password"
		as := NewAuthService(&userRepo)
		_, err := as.Login(context.Background(), emailInput)
		require.Error(t, err)
		require.Equal(t, err, data.ErrBadCredentials)

		require.True(t, userRepo.AssertNotCalled(t, "GetByUsername"))
		require.True(t, userRepo.AssertExpectations(t))
	})

	t.Run("invalid input", func(t *testing.T) {
		var userRepo mocks.UserRepo
		as := NewAuthService(&userRepo)
		_, err := as.Login(context.Background(), data.LoginInput{})
		require.Error(t, err)

		require.True(t, userRepo.AssertNotCalled(t, "GetByUsername"))
		require.True(t, userRepo.AssertNotCalled(t, "GetByEmail"))
		require.True(t, userRepo.AssertExpectations(t))
	})
}
