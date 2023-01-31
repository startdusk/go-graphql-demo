package postgres

import (
	"context"

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
		SELECT * FROM tweets
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

func (tr *TweetRepo) GetByID(ctx context.Context, userID, tweetID string) (data.Tweet, error) {
	const q = `
		SELECT * FROM tweets WHERE id=$1 AND user_id=$2 LIMIT 1
	`
	var t data.Tweet
	err := pgxscan.Get(ctx, tr.DB.Pool, &t, q, tweetID, userID)
	return t, err
}
