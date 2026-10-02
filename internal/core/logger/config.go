package core_logger

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит настройки для инициализации логгера.
type Config struct {
	Level  string `envconfig:"LEVEL" default:"info"`
	Folder string `envconfig:"FOLDER" default:""`
}

//!Сделал возвращаемое значение по указателю
// NewConfig считывает конфигурацию логгера из переменных окружения (префикс LOGGER_).
func NewConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return nil, fmt.Errorf("process envconfig: %w", err)
	}

	config.Level = strings.ToLower(strings.TrimSpace(config.Level))

	return &config, nil
}

// MustConfig считывает конфигурацию и вызывает panic при ошибке.
func MustConfig() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get logger config: %w", err))
	}
	return config
}