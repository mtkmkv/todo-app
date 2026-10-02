package users_transport_http

import (
	"net/http"

	core_logger "github.com/mtkmkv/todo-app/internal/core/logger"
	core_http_response "github.com/mtkmkv/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/mtkmkv/todo-app/internal/core/transport/http/utils"
)

// GetUserResponse структура ответа для получения пользователя по идентификатору
type GetUserResponse UserDTOResponse

func (h *UserHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// 1. Извлекаем и валидируем userID из пути запроса
	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	// 2. Получаем пользователя из бизнес-логики через сервис
	user, err := h.UserService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	// 3. Отдаем успешный результат со статусом 200 OK
	response := GetUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}