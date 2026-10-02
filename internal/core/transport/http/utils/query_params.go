package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	core_errors "github.com/mtkmkv/todo-app/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	query := r.URL.Query()
	if !query.Has(key) {
		return nil, nil
	}

	param := strings.TrimSpace(query.Get(key))
	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param '%s' by key '%s' is not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}
