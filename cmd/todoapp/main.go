package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/mtkmkv/todo-app/internal/core/logger"
	core_postgres_pool "github.com/mtkmkv/todo-app/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/mtkmkv/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/mtkmkv/todo-app/internal/core/transport/http/server"
	users_postgres_repository "github.com/mtkmkv/todo-app/internal/features/users/repository/postgres"
	users_service "github.com/mtkmkv/todo-app/internal/features/users/service"
	users_transport_http "github.com/mtkmkv/todo-app/internal/features/users/transport/http"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(*core_logger.MustConfig())
	if err != nil {
		fmt.Println("failed to initialize logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("init postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, *core_postgres_pool.MustConfig())
	if err != nil {
		logger.Error("failed to init postgres connection pool", zap.Error(err))
		return
	}
	defer pool.Close()

	logger.Debug("init feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	userTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

	logger.Debug("init HTTP Server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.MustConfig(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersionV1)
	apiVersionRouter.RegisterRoutes(userTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
