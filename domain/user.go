package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/startdusk/twitter/data"
)

type UserService struct {
	userRepo data.UserRepo
}

func NewUserService(ur data.UserRepo) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

func (us *UserService) GetByID(ctx context.Context, id string) (data.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return data.NilUser, data.ErrInvalidUUID
	}
	return us.userRepo.GetByID(ctx, id)
}
