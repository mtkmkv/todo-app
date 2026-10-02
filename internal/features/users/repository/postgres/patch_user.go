package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mtkmkv/todo-app/internal/core/domain"
	core_errors "github.com/mtkmkv/todo-app/internal/core/errors"
)

func (r *UserRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `
		UPDATE todoapp.users
		SET
			full_name = $1,
			phone_number = $2,
			version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING id, version, full_name, phone_number;
	`

	var model UserModel
	err := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, id, user.Version).Scan(
		&model.ID,
		&model.Version,
		&model.FullName,
		&model.PhoneNumber,
	)

	if err != nil {
		// Ошибка конкурентного изменения (версия не совпала) или пользователь был удален
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d' concurrently modified or deleted: %w", id, core_errors.ErrConflict)
		}

		// Ошибка дубликата номера телефона
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
			return domain.User{}, fmt.Errorf("%w: user with this phone already exists", core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("execute patch user query: %w", err)
	}

	return model.ToDomain(), nil
}
