package core_postgres_pool

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит параметры для подключения к пулу базы данных PostgreSQL.
type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

// NewConfig считывает конфигурацию пула PostgreSQL из переменных окружения (префикс POSTGRES_).
func NewConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return nil, fmt.Errorf("process envconfig: %w", err)
	}

	config.Host = strings.TrimSpace(config.Host)
	config.Port = strings.TrimSpace(config.Port)
	config.User = strings.TrimSpace(config.User)
	config.Database = strings.TrimSpace(config.Database)

	return &config, nil
}

// MustConfig считывает конфигурацию и вызывает panic при ошибке.
func MustConfig() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get postgres pool config: %w", err))
	}
	return config
}