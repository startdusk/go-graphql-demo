package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/startdusk/twitter/data"
)

type UserRepo struct {
	DB *DB
}

func (ur *UserRepo) Create(ctx context.Context, user data.User) (data.User, error) {
	const q = `
		INSERT INTO users (username, email, password) VALUES ($1, $2, $3)
		RETURNING *;
	`
	var newUser data.User
	err := pgxscan.Get(ctx, ur.DB.Pool, &newUser, q, user.Username, user.Email, user.Password)

	return newUser, err
}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (data.User, error) {
	const q = `
		SELECT * FROM users WHERE username = $1 LIMIT 1;
	`

	var user data.User
	if err := pgxscan.Get(ctx, ur.DB.Pool, &user, q, username); err != nil {
		if pgxscan.NotFound(err) {
			return data.NilUser, data.ErrNotFound
		}
		return data.NilUser, fmt.Errorf("select user by username: %s error: %w", username, err)
	}
	return user, nil
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (data.User, error) {
	const q = `
		SELECT * FROM users WHERE email = $1 LIMIT 1;
	`

	var user data.User
	if err := pgxscan.Get(ctx, ur.DB.Pool, &user, q, email); err != nil {
		if pgxscan.NotFound(err) {
			return data.NilUser, data.ErrNotFound
		}
		return data.NilUser, fmt.Errorf("select user by email: %s error: %w", email, err)
	}
	return user, nil
}

func (ur *UserRepo) GetByID(ctx context.Context, id string) (data.User, error) {
	const q = `
		SELECT * FROM users WHERE id = $1 LIMIT 1;
	`

	var user data.User
	if err := pgxscan.Get(ctx, ur.DB.Pool, &user, q, id); err != nil {
		if pgxscan.NotFound(err) {
			return data.NilUser, data.ErrNotFound
		}
		return data.NilUser, fmt.Errorf("select user by id: %s error: %w", id, err)
	}
	return user, nil
}

func (ur *UserRepo) GetByIDs(ctx context.Context, ids []string) ([]data.User, error) {
	const q = `
		SELECT * FROM users WHERE id = ANY($1)
	`

	var us []data.User
	err := pgxscan.Select(ctx, ur.DB.Pool, &us, q, ids)
	return us, err
}
