package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mtkmkv/todo-app/internal/core/domain"
	core_errors "github.com/mtkmkv/todo-app/internal/core/errors"
)

func (r *UserRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `
		SELECT id, version, full_name, phone_number
		FROM todoapp.users
		WHERE id = $1;
	`

	var model UserModel
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&model.ID,
		&model.Version,
		&model.FullName,
		&model.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d' not found: %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("execute get user query: %w", err)
	}

	return model.ToDomain(), nil
}