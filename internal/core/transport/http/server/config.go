package core_http_server

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит параметры для настройки и запуска HTTP-сервера.
type Config struct {
	Addr            string        `envconfig:"ADDR" default:":8080"`
	ShutDownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}

// NewConfig считывает конфигурацию HTTP-сервера из переменных окружения (префикс HTTP_).
func NewConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return nil, fmt.Errorf("process envconfig: %w", err)
	}

	config.Addr = strings.TrimSpace(config.Addr)

	return &config, nil
}

// MustConfig считывает конфигурацию и вызывает panic при ошибке.
func MustConfig() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get HTTP server config: %w", err))
	}
	return config
}