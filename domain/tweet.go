package domain

import (
	"context"

	"github.com/google/uuid"
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
	return ts.tweetRepo.All(ctx)
}

func (ts *TweetService) Create(ctx context.Context, input data.CreateTweetInput) (data.Tweet, error) {
	userID, err := shared.GetUserIDFromContext(ctx)
	if err != nil {
		return data.NilTweet, data.ErrUnauthenticated
	}
	input.Sanitize()
	if err := input.Validate(); err != nil {
		return data.NilTweet, err
	}

	return ts.tweetRepo.Create(ctx, data.Tweet{Body: input.Body, UserID: userID})
}

func (ts *TweetService) GetByID(ctx context.Context, tweetID string) (data.Tweet, error) {
	if _, err := uuid.Parse(tweetID); err != nil {
		return data.NilTweet, data.ErrInvalidUUID
	}

	return ts.tweetRepo.GetByID(ctx, tweetID)
}

func (ts *TweetService) Delete(ctx context.Context, tweetID string) error {
	curUserID, err := shared.GetUserIDFromContext(ctx)
	if err != nil {
		return data.ErrUnauthenticated
	}

	if _, err := uuid.Parse(tweetID); err != nil {
		return data.ErrInvalidUUID
	}

	tweet, err := ts.GetByID(ctx, tweetID)
	if err != nil {
		return err
	}

	if !tweet.CanDelete(curUserID) {
		return data.ErrForbidden
	}

	return ts.tweetRepo.Delete(ctx, tweetID)
}

func (ts *TweetService) CreateReply(ctx context.Context, parentID string, input data.CreateTweetInput) (data.Tweet, error) {
	userID, err := shared.GetUserIDFromContext(ctx)
	if err != nil {
		return data.NilTweet, data.ErrUnauthenticated
	}
	input.Sanitize()
	if err := input.Validate(); err != nil {
		return data.NilTweet, err
	}
	if _, err := uuid.Parse(parentID); err != nil {
		return data.NilTweet, data.ErrInvalidUUID
	}

	if _, err := ts.GetByID(ctx, parentID); err != nil {
		return data.NilTweet, data.ErrNotFound
	}

	return ts.tweetRepo.Create(ctx, data.Tweet{Body: input.Body, UserID: userID, ParentID: &parentID})
}
