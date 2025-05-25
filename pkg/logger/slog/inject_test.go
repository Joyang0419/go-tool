package slogx

import (
	"context"
	"log/slog"
	"testing"

	"go-tool/pkg/logger"
	"go-tool/pkg/logger/consts"
)

func TestFXModule(t *testing.T) {
	FXModule(Config{})

	ctx := context.WithValue(context.Background(), consts.TraceIDKey, "test-trace-id")
	logger.Log.InfoContext(ctx, "test slogx", slog.String("key", "value"))
}
