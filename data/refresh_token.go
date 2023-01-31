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
	ID         string
	Name       string
	UserID     string
	LastUsedAt time.Time
	ExpiredAt  time.Time
	CreatedAt  time.Time
}

type CreateRefreshTokenParams struct {
	Sub  string
	Name string
}

type RefreshTokenRepo interface {
	Create(ctx context.Context, params CreateRefreshTokenParams) (RefreshToken, error)
	GetByID(ctx context.Context, id string) (RefreshToken, error)
}
