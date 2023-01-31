package data

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTweetInput_Sanitize(t *testing.T) {
	input := CreateTweetInput{
		Body: "      body",
	}
	want := CreateTweetInput{
		Body: "body",
	}
	input.Sanitize()
	assert.Equal(t, input, want)
}

func TestCreateTweetInput_Validate(t *testing.T) {
	cases := []struct {
		name  string
		input CreateTweetInput
		err   error
	}{
		{
			name:  "valid",
			input: CreateTweetInput{Body: "hello"},
		},
		{
			name:  "tweet not long enough",
			input: CreateTweetInput{Body: "h"},
			err:   ErrValidation,
		},
		{
			name:  "tweet not long enough",
			input: CreateTweetInput{Body: strings.Repeat("s", 241)},
			err:   ErrValidation,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			err := cc.input.Validate()
			require.ErrorIs(t, err, cc.err)
		})
	}
}
