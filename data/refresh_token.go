package data

import (
	"context"
	"time"
)

var (
	AccessTokenLifetime  = 15 * time.Minute   // 15 minutes
	RefreshTokenLifetime = 7 * 24 * time.Hour // 1 week
)

type RefreshToken struct {
	ID         int64
	TokenID    string
	UserID     string
	LastUsedAt time.Time
	CreatedAt  time.Time
	ExpiredAt  time.Time
}

type CreateRefreshTokenParams struct {
	Sub       string
	TokenID   string
	ExpiredAt time.Time
}

type RefreshTokenRepo interface {
	Create(ctx context.Context, params CreateRefreshTokenParams) (RefreshToken, error)
	LastUsed(ctx context.Context, params CreateRefreshTokenParams) error
	GetByTokenID(ctx context.Context, tokenID string) (RefreshToken, error)
	Delete(ctx context.Context, tokenID string) error
}
