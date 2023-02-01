package tests

import (
	"context"
	"testing"

	"github.com/startdusk/twitter/data"
	"github.com/startdusk/twitter/shared"
	"github.com/stretchr/testify/assert"
)

func TestTweetService_Lifecycle(t *testing.T) {
	ctx := context.Background()
	resp, err := authService.Register(ctx, data.RegisterInput{
		Username:        "TestTweetService_Lifecycle",
		Email:           "TestTweetService_Lifecycle@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	})
	assert.NoError(t, err)

	ctx = context.WithValue(ctx, shared.UserIDKey{}, resp.User.ID)
	input := data.CreateTweetInput{
		Body: "hello",
	}
	initTweet, err := tweetService.Create(ctx, input)
	assert.NoError(t, err)

	cases := []struct {
		name    string
		op      func() (data.Tweet, error)
		wantErr bool
		err     error
	}{
		{
			name: "can create a tweet 1",
			op: func() (data.Tweet, error) {
				input := data.CreateTweetInput{
					Body: "hello",
				}
				return tweetService.Create(ctx, input)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "can create a tweet 2",
			op: func() (data.Tweet, error) {
				input := data.CreateTweetInput{
					Body: "hello",
				}
				return tweetService.Create(ctx, input)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "no auth user cannot create a tweet",
			op: func() (data.Tweet, error) {
				input := data.CreateTweetInput{
					Body: "hello",
				}
				return tweetService.Create(context.Background(), input)
			},
			wantErr: true,
			err:     data.ErrUnauthenticated,
		},
		{
			name: "tweet too short",
			op: func() (data.Tweet, error) {
				input := data.CreateTweetInput{
					Body: "h",
				}
				return tweetService.Create(ctx, input)
			},
			wantErr: true,
			err:     data.ErrValidation,
		},
		{
			name: "invalid tweet id",
			op: func() (data.Tweet, error) {
				return tweetService.GetByID(ctx, "xxx")
			},
			wantErr: true,
			err:     data.ErrInvalidUUID,
		},
		{
			name: "can get tweet by id",
			op: func() (data.Tweet, error) {
				return tweetService.GetByID(ctx, initTweet.ID)
			},
		},
		{
			name: "cannot delete tweet",
			op: func() (data.Tweet, error) {
				ctx := context.WithValue(context.Background(), shared.UserIDKey{}, "bad_user_id")
				return data.NilTweet, tweetService.Delete(ctx, initTweet.ID)
			},
			wantErr: true,
			err:     data.ErrForbidden,
		},
		{
			name: "can delete tweet",
			op: func() (data.Tweet, error) {
				return data.NilTweet, tweetService.Delete(ctx, initTweet.ID)
			},
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			tweet, err := cc.op()
			if cc.wantErr {
				assert.Error(t, err)
				if cc.err != nil {
					assert.ErrorIs(t, err, cc.err)
				}
			} else {
				assert.NotNil(t, tweet)
			}
		})
	}

	tweets, err := tweetService.All(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(tweets), 2)
}
