package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/mtkmkv/todo-app/internal/core/logger"
	core_http_response "github.com/mtkmkv/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/mtkmkv/todo-app/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UserHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// 1. Извлекаем и валидируем query-параметры пагинации
	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit'/'offset' query param",
		)
		return
	}

	// 2. Получаем список пользователей через бизнес-логику сервиса
	userDomains, err := h.UserService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)
		return
	}

	// 3. Преобразуем доменные сущности в DTO и отдаем ответ со статусом 200 OK
	dtos := usersDTOFromDomains(userDomains)
	if dtos == nil {
		dtos = make([]UserDTOResponse, 0)
	}

	response := GetUsersResponse(dtos)
	responseHandler.JSONResponse(response, http.StatusOK)
}


// getLimitOffsetQueryParams извлекает параметры limit и offset из строки запроса
func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}