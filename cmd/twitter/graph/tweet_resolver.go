package graph

import (
	"context"
	"errors"

	"github.com/startdusk/twitter/data"
)

func (q *queryResolver) Tweets(ctx context.Context) ([]*Tweet, error) {
	ts, err := q.TweetService.All(ctx)
	if err != nil {
		return nil, err
	}

	tweets := make([]*Tweet, len(ts))
	for i, t := range ts {
		tweets[i] = mapToTweet(t)
	}

	return tweets, nil
}

func (m *mutationResolver) CreateTweet(ctx context.Context, input CreatedTweetInput) (*Tweet, error) {
	t, err := m.TweetService.Create(ctx, data.CreateTweetInput{
		Body: input.Body,
	})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUnauthenticated):
			return nil, writeUnauthenticatedError(ctx, err)
		default:
			return nil, err
		}
	}
	return mapToTweet(t), err
}

func mapToTweet(t data.Tweet) *Tweet {
	return &Tweet{
		ID:        t.ID,
		Body:      t.Body,
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
	}
}
