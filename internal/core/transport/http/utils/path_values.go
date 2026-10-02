package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	core_errors "github.com/mtkmkv/todo-app/internal/core/errors"
)

// GetIntPathValue извлекает положительный целочисленный параметр из пути запроса (например, /users/{id}).
func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := strings.TrimSpace(r.PathValue(key))
	if pathValue == "" {
		return 0, fmt.Errorf("no key '%s' in path value: %w", key, core_errors.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf(
			"path value='%s' by key='%s' is not a valid integer: %v: %w",
			pathValue,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	if val <= 0 {
		return 0, fmt.Errorf(
			"path value='%d' by key='%s' must be positive: %w",
			val,
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	return val, nil
}