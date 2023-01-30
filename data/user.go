package data

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUsernameTaken = errors.New("username has taken")
	ErrEmailTaken    = errors.New("email has taken")
)

var NilUser User

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepo interface {
	Create(ctx context.Context, user User) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}
