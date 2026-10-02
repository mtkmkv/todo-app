package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mtkmkv/todo-app/internal/core/domain"
	core_logger "github.com/mtkmkv/todo-app/internal/core/logger"
	core_http_request "github.com/mtkmkv/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/mtkmkv/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/mtkmkv/todo-app/internal/core/transport/http/types"
	core_http_utils "github.com/mtkmkv/todo-app/internal/core/transport/http/utils"
)

// PatchUserRequest структура входящего запроса для частичного обновления пользователя
type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 characters")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 characters")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must start with '+'")
			}
		}
	}

	return nil
}

// PatchUserResponse структура ответа после частичного обновления пользователя
type PatchUserResponse UserDTOResponse

func (h *UserHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// 1. Извлекаем и валидируем userID из пути запроса
	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	// 2. Декодируем и валидируем тело запроса
	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	// 3. Преобразуем транспортную DTO-модель патча в доменную
	userPatch := userPatchFromRequest(request)

	// 4. Применяем частичное обновление через сервис бизнес-логики
	userDomain, err := h.UserService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	// 5. Отдаем успешный результат с обновленным пользователем и статусом 200 OK
	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
