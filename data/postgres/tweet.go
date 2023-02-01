package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/startdusk/twitter/data"
)

type TweetRepo struct {
	DB *DB
}

func NewTweetRepo(db *DB) *TweetRepo {
	return &TweetRepo{DB: db}
}

func (tr *TweetRepo) All(ctx context.Context) ([]data.Tweet, error) {
	const q = `
		SELECT * FROM tweets ORDER BY created_at DESC
	`
	var ts []data.Tweet
	err := pgxscan.Select(ctx, tr.DB.Pool, &ts, q)
	return ts, err
}

func (tr *TweetRepo) Create(ctx context.Context, tweet data.Tweet) (data.Tweet, error) {
	const q = `
		INSERT INTO tweets (body, user_id) VALUES ($1, $2) RETURNING *
	`
	var t data.Tweet
	err := pgxscan.Get(ctx, tr.DB.Pool, &t, q, tweet.Body, tweet.UserID)
	return t, err
}

func (tr *TweetRepo) GetByID(ctx context.Context, tweetID string) (data.Tweet, error) {
	const q = `
		SELECT * FROM tweets WHERE id = $1 LIMIT 1
	`
	var t data.Tweet
	if err := pgxscan.Get(ctx, tr.DB.Pool, &t, q, tweetID); err != nil {
		if pgxscan.NotFound(err) {
			return data.NilTweet, data.ErrNotFound
		}
		return data.NilTweet, fmt.Errorf("select tweets by id=%s error: %w", tweetID, err)
	}
	return t, nil
}
