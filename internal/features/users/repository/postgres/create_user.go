package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mtkmkv/todo-app/internal/core/domain"
	core_errors "github.com/mtkmkv/todo-app/internal/core/errors"
)

const pgErrUniqueViolation = "23505"

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `
		INSERT INTO todoapp.users(full_name, phone_number)
		VALUES ($1, $2)
		RETURNING id, version, full_name, phone_number;
	`

	var model UserModel
	err := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber).Scan(
		&model.ID,
		&model.Version,
		&model.FullName,
		&model.PhoneNumber,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
			return domain.User{}, fmt.Errorf("%w: user with this phone already exists", core_errors.ErrConflict)
		}
		return domain.User{}, fmt.Errorf("execute create user query: %w", err)
	}

	return model.ToDomain(), nil
}
