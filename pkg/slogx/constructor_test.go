package slogx

import (
	"context"
	"log/slog"
	"testing"
)

func TestNewSlog(t *testing.T) {
	logger := NewSlog(Config{
		Path:       "",
		MaxSize:    0,
		MaxAge:     0,
		CallerSkip: 7,
	})

	slog.SetDefault(logger)

	slog.Info("test slogx", slog.String("key", "value"))
}

func TestNewLogCtx(t *testing.T) {
	logger := NewSlog(Config{
		Path:       "",
		MaxSize:    0,
		MaxAge:     0,
		CallerSkip: 7,
	})

	slog.SetDefault(logger)

	ctx := context.WithValue(context.Background(), TraceIDKey, "test-trace-id")

	slog.InfoContext(ctx, "test slogx")
}
