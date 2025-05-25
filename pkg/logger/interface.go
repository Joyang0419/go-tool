package logger

import (
	"context"
)

type ILogger interface {
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}
