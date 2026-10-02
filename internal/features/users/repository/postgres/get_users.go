package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/mtkmkv/todo-app/internal/core/domain"
)

func (r *UserRepository) GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `
		SELECT id, version, full_name, phone_number
		FROM todoapp.users
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	capacity := 0
	if limit != nil && *limit > 0 {
		capacity = *limit
	}
	userDomains := make([]domain.User, 0, capacity)

	for rows.Next() {
		var model UserModel
		err := rows.Scan(
			&model.ID,
			&model.Version,
			&model.FullName,
			&model.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}

		userDomains = append(userDomains, model.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user rows: %w", err)
	}

	return userDomains, nil
}