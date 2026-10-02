package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/mtkmkv/todo-app/internal/core/logger"
	core_http_response "github.com/mtkmkv/todo-app/internal/core/transport/http/response"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-ID"
)

// RequestID — middleware для генерации и проставления сквозного идентификатора запроса
// исключительно на уровне HTTP-заголовков
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// Logger — middleware, которая обогащает базовый логгер контекстом текущего запроса
// (request_id, URL) и сохраняет этот "дочерний" логгер в context.Context.
// Это позволяет всем последующим слоям приложения (контроллерам, сервисам) вести
// логирование с автоматической привязкой к конкретному HTTP-запросу.
func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			// Кладём логгер через функцию самого пакета core_logger
			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Panic — это middleware перехвата непредвиденных критических сбоев (паник).
// Она гарантирует, что даже если в контроллере или сервисе произойдет panic(),
// веб-сервер не упадет, а залогирует ошибку и корректно вернет клиенту 500 статус.
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			// Передаем w в конструктор
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					// Вызываем без передачи w
					responseHandler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Trace — middleware для сквозного логирования жизненного цикла HTTP-запроса.
// Она фиксирует время начала обработки, логирует факт входящего запроса,
// а по завершении (даже в случае паники) высчитывает латенцию (время обработки)
// и логирует итоговый статус-код ответа.
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			rw := core_http_response.NewResponseWriter(w)
			before := time.Now()

			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			defer func() {
				log.Debug(
					"<<< done HTTP request",
					zap.Int("status_code", rw.GetStatusCodeOrPanic()),
					zap.Duration("latency", time.Since(before)),
				)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}