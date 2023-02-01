package domain

import (
	"context"
	"testing"

	"github.com/startdusk/twitter/data"
	mocks "github.com/startdusk/twitter/mocks/data"
	"github.com/startdusk/twitter/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTweetService_All(t *testing.T) {
	t.Run("can get all tweets", func(t *testing.T) {
		var tweetRepo mocks.TweetRepo
		returnTweets := []data.Tweet{{
			ID:     "171046a0-a81f-473f-b8a5-d0f8e9e832be",
			Body:   "hello",
			UserID: "user_id",
		}}
		tweetRepo.On("All", mock.Anything, mock.Anything).Return(returnTweets, nil)
		tr := NewTweetService(&tweetRepo)
		ctx := context.Background()
		tweet, err := tr.All(ctx)
		assert.NoError(t, err)
		assert.Equal(t, tweet, returnTweets)
	})
}

func TestTweetService_Create(t *testing.T) {
	t.Run("can create tweet", func(t *testing.T) {
		var tweetRepo mocks.TweetRepo
		returnTweet := data.Tweet{
			ID:     "44bad1d3-946a-4869-aed0-bdc5e3e14c93",
			Body:   "hello",
			UserID: "user_id",
		}
		tweetRepo.On("Create", mock.Anything, mock.Anything).Return(returnTweet, nil)
		tr := NewTweetService(&tweetRepo)

		ctx := context.WithValue(context.Background(), shared.UserIDKey{}, "user_id")
		tweet, err := tr.Create(ctx, data.CreateTweetInput{
			Body: "hello",
		})
		assert.NoError(t, err)
		assert.Equal(t, tweet, returnTweet)
	})

	t.Run("not auth user cannot create a tweet", func(t *testing.T) {
		var tweetRepo mocks.TweetRepo
		tr := NewTweetService(&tweetRepo)

		ctx := context.Background()
		tweet, err := tr.Create(ctx, data.CreateTweetInput{
			Body: "hello",
		})
		assert.ErrorIs(t, err, data.ErrUnauthenticated)
		assert.Equal(t, tweet, data.NilTweet)

		tweetRepo.AssertNotCalled(t, "Create")
		tweetRepo.AssertExpectations(t)
	})
}

func TestTweetService_GetByID(t *testing.T) {
	t.Run("can get a tweet by id", func(t *testing.T) {
		var tweetRepo mocks.TweetRepo
		tweetID := "3d73c595-b771-4be1-9328-82e1cb93b350"
		returnTweet := data.Tweet{
			ID:     tweetID,
			Body:   "hello",
			UserID: "user_id",
		}
		tweetRepo.On("GetByID", mock.Anything, mock.Anything).Return(returnTweet, nil)
		tr := NewTweetService(&tweetRepo)

		ctx := context.WithValue(context.Background(), shared.UserIDKey{}, "user_id")
		tweet, err := tr.GetByID(ctx, tweetID)
		assert.NoError(t, err)
		assert.Equal(t, tweet, returnTweet)
	})

	t.Run("invalid tweet id", func(t *testing.T) {
		var tweetRepo mocks.TweetRepo
		tweetID := "xxx"
		tr := NewTweetService(&tweetRepo)

		ctx := context.Background()
		tweet, err := tr.GetByID(ctx, tweetID)
		assert.ErrorIs(t, err, data.ErrInvalidUUID)
		assert.Equal(t, tweet, data.NilTweet)

		tweetRepo.AssertNotCalled(t, "GetByID")
		tweetRepo.AssertExpectations(t)
	})
}

func TestTweetService_Delete(t *testing.T) {
	t.Run("can delete tweet", func(t *testing.T) {
		var tweetRepo mocks.TweetRepo
		tweetID := "3d73c595-b771-4be1-9328-82e1cb93b350"
		userID := "user_id"
		returnTweet := data.Tweet{
			ID:     tweetID,
			Body:   "hello",
			UserID: userID,
		}
		tweetRepo.On("GetByID", mock.Anything, mock.Anything).Return(returnTweet, nil)
		tweetRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)
		tr := NewTweetService(&tweetRepo)

		ctx := context.WithValue(context.Background(), shared.UserIDKey{}, userID)
		err := tr.Delete(ctx, tweetID)
		assert.NoError(t, err)
	})

	t.Run("invalid tweet id", func(t *testing.T) {
		var tweetRepo mocks.TweetRepo
		tweetID := "xxx"
		tr := NewTweetService(&tweetRepo)

		ctx := context.WithValue(context.Background(), shared.UserIDKey{}, "user_id")
		err := tr.Delete(ctx, tweetID)
		assert.ErrorIs(t, err, data.ErrInvalidUUID)

		tweetRepo.AssertNotCalled(t, "Delete")
		tweetRepo.AssertExpectations(t)
	})
}
