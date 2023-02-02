package postgres

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/startdusk/twitter/data"
)

type RefreshTokenRepo struct {
	DB *DB
}

func (usr *RefreshTokenRepo) Create(ctx context.Context, params data.CreateRefreshTokenParams) (data.RefreshToken, error) {
	const q = `
		INSERT INTO user_sessions (user_id, token_id, expired_at) VALUES ($1, $2, $3) RETURNING *
	`
	var token data.RefreshToken
	err := pgxscan.Get(ctx, usr.DB.Pool, &token, q, params.Sub, params.TokenID, params.ExpiredAt)
	return token, err
}

func (usr *RefreshTokenRepo) LastUsed(ctx context.Context, params data.CreateRefreshTokenParams) error {
	const q = `
		UPDATE user_sessions SET last_used_at = $1 WHERE token_id = $2 AND user_id = $3
	`

	_, err := usr.DB.Pool.Exec(ctx, q, time.Now(), params.TokenID, params.Sub)
	return err
}

func (usr *RefreshTokenRepo) Delete(ctx context.Context, tokenID string) error {
	const q = `
		DELETE FROM user_sessions WHERE token_id = $1
	`

	_, err := usr.DB.Pool.Exec(ctx, q, tokenID)
	return err
}

func (usr *RefreshTokenRepo) GetByTokenID(ctx context.Context, tokenID string) (data.RefreshToken, error) {
	const q = `
		SELECT * FROM user_sessions WHERE token_id = $1
	`
	var token data.RefreshToken
	err := pgxscan.Get(ctx, usr.DB.Pool, &token, q, tokenID)
	return token, err
}
