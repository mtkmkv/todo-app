package users_transport_http

import (
	"net/http"

	core_logger "github.com/mtkmkv/todo-app/internal/core/logger"
	core_http_response "github.com/mtkmkv/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/mtkmkv/todo-app/internal/core/transport/http/utils"
)

func (h *UserHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// 1. Извлекаем и валидируем userID из пути запроса
	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	// 2. Вызываем удаление пользователя через сервис
	if err := h.UserService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	// 3. Отдаем успешный результат со статусом 204 No Content
	responseHandler.NoContentResponse()
}