package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Неэкспортируемый уникальный тип для ключа контекста
type ctxKeyLogger struct{}

var loggerCtxKey = ctxKeyLogger{}

type Logger struct {
	*zap.Logger

	file *os.File
}

// ToContext сохраняет логгер в context.Context по единому приватному ключу.
func ToContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, log)
}

// FromContext безопасно извлекает логгер из context.Context.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(loggerCtxKey).(*Logger)
	if !ok {
		panic("logger not found in context")
	}
	return log
}

// NewLogger инициализирует и настраивает новый экземпляр логгера.
func NewLogger(config Config) (*Logger, error) {
	// Установка и парсинг уровня логирования (например: debug, info, warn, error)
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	// Автоматическое создание цепочки папок для логов, если их еще нет
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("create log folder: %w", err)
	}

	// Формирование имени файла с таймстемпом
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	// Настройка конфигурации энкодера отображения логов
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	// Создание Tee-ядра для логирования параллельно в stdout и файл
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

// With создает дочерний экземпляр логгера с добавлением постоянных полей контекста.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}

// Close сбрасывает все накопившиеся буферы логов и корректно закрывает дескриптор файла.
func (l *Logger) Close() {
	_ = l.Logger.Sync()

	if l.file != nil {
		if err := l.file.Close(); err != nil {
			fmt.Println("Error closing log file:", err)
		}
	}
}