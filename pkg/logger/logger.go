package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"tgVideoCall/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func MustInitLogger(cfg config.Config) slog.Logger {
	var logFile *os.File
	var err error

	if cfg.Log.FilePath != "" { //Если строка в конфиге пустая, это будет означать что нам не нужно сохранение логов в файл
		logFile, err = os.OpenFile(cfg.Log.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("error opening file:", err)
		}
	}

	var log *slog.Logger

	switch cfg.Env {
	case envLocal:
		if cfg.Log.FilePath == "" {
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			return *log
		}
		log = slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		if cfg.Log.FilePath == "" {
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			return *log
		}
		log = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	if cfg.Log.FilePath != "" {
		log.Info(fmt.Sprintf("Logs are saving to: %s", cfg.Log.FilePath))
	}

	return *log
}

// ZapSlogHandler реализует slog.Handler для zap.Logger
type ZapSlogHandler struct {
	zapLogger *zap.Logger
}

func NewZapSlogHandler(zapLogger *zap.Logger) *ZapSlogHandler {
	return &ZapSlogHandler{zapLogger: zapLogger}
}

func (h *ZapSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Преобразуем slog.Level в zapcore.Level
	var zapLevel zapcore.Level
	switch {
	case level < slog.LevelDebug:
		zapLevel = zapcore.DebugLevel
	case level < slog.LevelInfo:
		zapLevel = zapcore.InfoLevel
	case level < slog.LevelWarn:
		zapLevel = zapcore.WarnLevel
	case level < slog.LevelError:
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.DPanicLevel
	}
	return h.zapLogger.Core().Enabled(zapLevel)
}

func (h *ZapSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	// Преобразуем slog.Record в поля zap
	fields := make([]zap.Field, 0, r.NumAttrs()+1)
	fields = append(fields, zap.String("message", r.Message))

	// Добавляем атрибуты
	r.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, zap.Any(attr.Key, attr.Value.Any()))
		return true
	})

	// Добавляем источник (source), если есть
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		fields = append(fields, zap.String("source", f.File+":"+strconv.FormatUint(uint64(r.PC), 10)))
	}

	// Логируем с соответствующим уровнем
	switch r.Level {
	case slog.LevelDebug:
		h.zapLogger.Debug(r.Message, fields...)
	case slog.LevelInfo:
		h.zapLogger.Info(r.Message, fields...)
	case slog.LevelWarn:
		h.zapLogger.Warn(r.Message, fields...)
	case slog.LevelError:
		h.zapLogger.Error(r.Message, fields...)
	default:
		h.zapLogger.DPanic(r.Message, fields...)
	}

	return nil
}

func (h *ZapSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]zap.Field, len(attrs))
	for i, attr := range attrs {
		fields[i] = zap.Any(attr.Key, attr.Value.Any())
	}
	return &ZapSlogHandler{zapLogger: h.zapLogger.With(fields...)}
}

func (h *ZapSlogHandler) WithGroup(name string) slog.Handler {
	// Для простоты игнорируем группы, но в production нужно реализовать
	return h
}

// WrapZapToSlog оборачивает zap.Logger в slog.Logger
func WrapZapToSlog(zapLogger *zap.Logger) *slog.Logger {
	return slog.New(NewZapSlogHandler(zapLogger))
}
