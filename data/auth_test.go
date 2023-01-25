package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterInput_Sanitize(t *testing.T) {
	actual := RegisterInput{
		Username:        "            bob",
		Email:           "bob@gmail.com              ",
		Password:        "password         ",
		ConfirmPassword: "            password         ",
	}

	expected := RegisterInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	actual.Sanitize()
	require.Equal(t, expected, actual)
}

func TestRegisterInput_Validate(t *testing.T) {
	cases := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name: "valid",
			input: RegisterInput{
				Username:        "bob",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: nil,
		},
		{
			name: "username_too_short",
			input: RegisterInput{
				Username:        "b",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "password_too_short",
			input: RegisterInput{
				Username:        "bob",
				Email:           "bob@gmail.com",
				Password:        "",
				ConfirmPassword: "",
			},
			err: ErrValidation,
		},
		{
			name: "email_invalid",
			input: RegisterInput{
				Username:        "bob",
				Email:           "bobxx.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "password_no_equal_confirm_password",
			input: RegisterInput{
				Username:        "bob",
				Email:           "bob@gmail.com",
				Password:        "password1",
				ConfirmPassword: "password2",
			},
			err: ErrValidation,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if err != nil {
				require.ErrorIs(t, err, ErrValidation)
			}
		})
	}
}

func TestIsEmail(t *testing.T) {
	require.True(t, IsEmail("123@gmail.com"))
	require.False(t, IsEmail("123gmail.com"))
}
