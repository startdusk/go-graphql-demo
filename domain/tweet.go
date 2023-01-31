package domain

import (
	"context"

	"github.com/startdusk/twitter/data"
	"github.com/startdusk/twitter/shared"
)

type TweetService struct {
	tweetRepo data.TweetRepo
}

func NewTweetService(tr data.TweetRepo) *TweetService {
	return &TweetService{tweetRepo: tr}
}

func (ts *TweetService) All(ctx context.Context) ([]data.Tweet, error) {
	_, err := shared.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, data.ErrUnauthenticated
	}
	return ts.tweetRepo.All(ctx)
}

func (ts *TweetService) Create(ctx context.Context, input data.CreateTweetInput) (data.Tweet, error) {
	_, err := shared.GetUserIDFromContext(ctx)
	if err != nil {
		return data.NilTweet, data.ErrUnauthenticated
	}
	return ts.tweetRepo.Create(ctx, input)
}

func (ts *TweetService) GetByID(ctx context.Context, id string) (data.Tweet, error) {
	_, err := shared.GetUserIDFromContext(ctx)
	if err != nil {
		return data.NilTweet, data.ErrUnauthenticated
	}
	return ts.tweetRepo.GetByID(ctx, id)
}
