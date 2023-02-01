package graph

import (
	"context"

	"github.com/startdusk/twitter/data"
	"github.com/startdusk/twitter/shared"
)

func (q *queryResolver) Me(ctx context.Context) (*User, error) {
	userID, err := shared.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, data.ErrUnauthenticated
	}

	return &User{ID: userID}, nil
}

func mapToUser(user data.User) *User {
	return &User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
