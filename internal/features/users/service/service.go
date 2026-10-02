package user_service

import (
	"context"

	"github.com/mtkmkv/todo-app/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

func NewUserService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
