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

func (m *mutationResolver) DeleteTweet(ctx context.Context, tweetID string) (bool, error) {
	err := m.TweetService.Delete(ctx, tweetID)
	return err == nil, err
}

func (m *mutationResolver) CreateReply(ctx context.Context, parentID string, input CreatedTweetInput) (*Tweet, error) {
	t, err := m.TweetService.CreateReply(ctx, parentID, data.CreateTweetInput{
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

func (t *tweetResolver) User(ctx context.Context, obj *Tweet) (*User, error) {
	// 使用loder 加载数据到内存, 减少sql查询次数
	return DataloaderFor(ctx).UserByID.Load(obj.UserID)
	// user, err := t.UserService.GetByID(ctx, obj.UserID)
	// if err != nil {
	// 	return nil, err
	// }
	// return mapToUser(user), err
}

func mapToTweet(t data.Tweet) *Tweet {
	return &Tweet{
		ID:        t.ID,
		Body:      t.Body,
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
	}
}
