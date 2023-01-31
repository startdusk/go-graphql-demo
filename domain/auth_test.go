package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/startdusk/twitter/data"
	"github.com/startdusk/twitter/faker"
	mocks "github.com/startdusk/twitter/mocks/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

		var authTokenService mocks.AuthTokenService
		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("a access token", nil)

		as := NewAuthService(&userRepo, &authTokenService)
		resp, err := as.Register(context.Background(), input)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.User.ID)
		assert.Equal(t, input.Email, resp.User.Email)
		assert.Equal(t, input.Username, resp.User.Username)
		assert.NotEmpty(t, resp.AccessToken)

		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
	})

	t.Run("cannot generate access token", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{
			ID:       "user_id",
			Username: "bob",
			Email:    "bob@gmail.com",
		}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(data.NilUser, nil)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(data.NilUser, nil)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(returnUser, nil)

		var authTokenService mocks.AuthTokenService
		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("", errors.New("error"))

		as := NewAuthService(&userRepo, &authTokenService)
		resp, err := as.Register(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, resp, data.NilAuthResponse)
		assert.Equal(t, err, data.ErrGenAccessToken)
	})

	t.Run("username taken", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(returnUser, data.ErrUsernameTaken)

		var authTokenService mocks.AuthTokenService

		as := NewAuthService(&userRepo, &authTokenService)
		_, err := as.Register(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, err, data.ErrUsernameTaken)

		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)

		authTokenService.AssertNotCalled(t, "CreateAccessToken")
		authTokenService.AssertExpectations(t)
	})

	t.Run("email taken", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(returnUser, nil)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, data.ErrEmailTaken)

		var authTokenService mocks.AuthTokenService

		as := NewAuthService(&userRepo, &authTokenService)
		_, err := as.Register(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, err, data.ErrEmailTaken)

		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)

		authTokenService.AssertNotCalled(t, "CreateAccessToken")
		authTokenService.AssertExpectations(t)
	})

	t.Run("create error", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(returnUser, nil)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, nil)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(returnUser, errors.New("create user error"))

		var authTokenService mocks.AuthTokenService

		as := NewAuthService(&userRepo, &authTokenService)
		_, err := as.Register(context.Background(), input)
		assert.Error(t, err)

		userRepo.AssertExpectations(t)

		authTokenService.AssertNotCalled(t, "CreateAccessToken")
		authTokenService.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		var userRepo mocks.UserRepo
		var authTokenService mocks.AuthTokenService

		as := NewAuthService(&userRepo, &authTokenService)
		_, err := as.Register(context.Background(), data.RegisterInput{})
		assert.Error(t, err)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)

		authTokenService.AssertNotCalled(t, "CreateAccessToken")
		authTokenService.AssertExpectations(t)
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

		var authTokenService mocks.AuthTokenService
		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("a access token", nil)

		as := NewAuthService(&userRepo, &authTokenService)
		resp, err := as.Login(context.Background(), usernameInput)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.User.ID)
		assert.Equal(t, usernameInput.UsernameOrEmail, resp.User.Username)
		assert.NotEmpty(t, resp.AccessToken)

		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertExpectations(t)
		authTokenService.AssertExpectations(t)
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

		var authTokenService mocks.AuthTokenService
		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("a access token", nil)

		as := NewAuthService(&userRepo, &authTokenService)
		resp, err := as.Login(context.Background(), emailInput)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.User.ID)
		assert.Equal(t, emailInput.UsernameOrEmail, resp.User.Email)
		assert.NotEmpty(t, resp.AccessToken)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertExpectations(t)

		authTokenService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, data.ErrNotFound)
		var authTokenService mocks.AuthTokenService

		as := NewAuthService(&userRepo, &authTokenService)
		_, err := as.Login(context.Background(), emailInput)
		assert.Error(t, err)
		assert.Equal(t, err, data.ErrBadCredentials)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertExpectations(t)

		authTokenService.AssertNotCalled(t, "CreateAccessToken")
		authTokenService.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		var userRepo mocks.UserRepo
		returnUser := data.User{}
		expectedErr := errors.New("internal error")
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(returnUser, expectedErr)
		var authTokenService mocks.AuthTokenService

		as := NewAuthService(&userRepo, &authTokenService)
		_, err := as.Login(context.Background(), emailInput)
		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertExpectations(t)

		authTokenService.AssertNotCalled(t, "CreateAccessToken")
		authTokenService.AssertExpectations(t)
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
		var authTokenService mocks.AuthTokenService

		as := NewAuthService(&userRepo, &authTokenService)
		_, err := as.Login(context.Background(), emailInput)
		assert.Error(t, err)
		assert.Equal(t, err, data.ErrBadCredentials)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertExpectations(t)

		authTokenService.AssertNotCalled(t, "CreateAccessToken")
		authTokenService.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		var userRepo mocks.UserRepo
		var authTokenService mocks.AuthTokenService

		as := NewAuthService(&userRepo, &authTokenService)
		_, err := as.Login(context.Background(), data.LoginInput{})
		assert.Error(t, err)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertExpectations(t)

		authTokenService.AssertNotCalled(t, "CreateAccessToken")
		authTokenService.AssertExpectations(t)
	})
}
