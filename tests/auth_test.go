package tests

import (
	"context"
	"testing"

	"github.com/startdusk/twitter/data"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Register_Lifecycle(t *testing.T) {
	cases := []struct {
		name    string
		input   data.RegisterInput
		wantErr bool
		err     error
	}{
		{
			name: "can register",
			input: data.RegisterInput{
				Username:        "bob",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
		},
		{
			name: "username taken",
			input: data.RegisterInput{
				Username:        "bob",
				Email:           "bob1@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			wantErr: true,
			err:     data.ErrUsernameTaken,
		},
		{
			name: "email taken",
			input: data.RegisterInput{
				Username:        "bob1",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			wantErr: true,
			err:     data.ErrEmailTaken,
		},
		{
			name:    "invalid input",
			input:   data.RegisterInput{},
			wantErr: true,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			resp, err := authService.Register(context.Background(), cc.input)
			if cc.wantErr {
				assert.Error(t, err)
				if cc.err != nil {
					assert.ErrorIs(t, err, cc.err)
				}
			} else {
				assert.NotEmpty(t, resp.User.ID)
				assert.Equal(t, cc.input.Email, resp.User.Email)
				assert.Equal(t, cc.input.Username, resp.User.Username)
				assert.NotEmpty(t, resp.AccessToken)
			}
		})
	}
}

func TestAuthService_Login_Lifecycle(t *testing.T) {
	username := "login_test"
	email := "login_test@gamil.com"
	password := "password"

	_, err := authService.Register(context.Background(), data.RegisterInput{
		Username:        username,
		Email:           email,
		Password:        password,
		ConfirmPassword: password,
	})
	if err != nil {
		panic(err)
	}

	cases := []struct {
		name    string
		input   data.LoginInput
		wantErr bool
		err     error
	}{
		{
			name: "username can login",
			input: data.LoginInput{
				UsernameOrEmail: "login_test",
				Password:        password,
			},
		},
		{
			name: "email can login",
			input: data.LoginInput{
				UsernameOrEmail: "login_test@gmail.com",
				Password:        password,
			},
		},
		{
			name: "user not found",
			input: data.LoginInput{
				UsernameOrEmail: "unkown user",
				Password:        password,
			},
			wantErr: true,
			err:     data.ErrBadCredentials,
		},
		{
			name: "user password error",
			input: data.LoginInput{
				UsernameOrEmail: "login_test",
				Password:        "invalid_password",
			},
			wantErr: true,
			err:     data.ErrBadCredentials,
		},
		{
			name:    "invalid input",
			input:   data.LoginInput{},
			wantErr: true,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			resp, err := authService.Login(context.Background(), cc.input)
			if err != nil {
				assert.Error(t, err)
				if cc.err != nil {
					assert.Equal(t, err, cc.err)
				}
			} else {
				assert.NotEmpty(t, resp.User.ID)
				assert.NotEmpty(t, resp.AccessToken)
			}
		})
	}
}
