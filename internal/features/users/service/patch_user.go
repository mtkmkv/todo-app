package user_service

import (
	"context"
	"fmt"

	"github.com/mtkmkv/todo-app/internal/core/domain"
)

func (s *UsersService) PatchUser(ctx context.Context, id int, patch domain.UserPatch) (domain.User, error) {
	// 1. Получаем текущее состояние пользователя из БД
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	// 2. Применяем патч к доменной модели (валидация бизнес-правил)
	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	// 3. Сохраняем обновленные данные в репозитории по целевому id
	patchedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user in repository: %w", err)
	}

	return patchedUser, nil
}