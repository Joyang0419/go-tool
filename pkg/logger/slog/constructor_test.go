package slog

import (
	"context"
	"log/slog"
	"testing"

	"go-tool/pkg/logger/consts"
)

func TestNewSlog(t *testing.T) {
	logger := New(Config{
		Path:    "",
		MaxSize: 0,
		MaxAge:  0,
	})

	slog.SetDefault(logger)

	slog.Info("test slogx", slog.String("key", "value"))
}

func TestNewLogCtx(t *testing.T) {
	l := New(
		Config{
			Path:    "",
			MaxSize: 0,
			MaxAge:  0,
		})

	slog.SetDefault(l)

	ctx := context.WithValue(context.Background(), consts.TraceIDKey, "test-trace-id")

	slog.InfoContext(ctx, "test slogx")
}
