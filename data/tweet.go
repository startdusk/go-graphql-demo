package data

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var (
	TweetMinLength = 2
	TweetMaxLength = 240
)

type CreateTweetInput struct {
	Body string
}

func (in *CreateTweetInput) Sanitize() {
	in.Body = strings.TrimSpace(in.Body)
}

func (in CreateTweetInput) Validate() error {
	if len(in.Body) < TweetMinLength {
		return fmt.Errorf("%w: body not long enough, (%d) characters at least", ErrValidation, TweetMinLength)
	}

	if len(in.Body) > TweetMaxLength {
		return fmt.Errorf("%w: body too long enough, (%d) characters at max", ErrValidation, TweetMaxLength)
	}
	return nil
}

var NilTweet Tweet

type Tweet struct {
	ID        string
	Body      string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TweetService interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, input CreateTweetInput) (Tweet, error)
	GetByID(ctx context.Context, id string) (Tweet, error)
}

type TweetRepo interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, tweet Tweet) (Tweet, error)
	GetByID(ctx context.Context, tweetID, userID string) (Tweet, error)
}
