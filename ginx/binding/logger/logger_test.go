package logger_test

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"

	"go-tool/ginx/binding/logger"
	"go-tool/ginx/consts"
)

func TestLogger_Info(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, consts.TraceIDKey, "test_trace_id")
	log := logger.NewDefaultLogger()
	log.Info(ctx, "test info message", zap.String("test_key", "test_value"))
}

func TestLogger_Error(t *testing.T) {
	ctx := context.Background()
	msg := "test error message"
	log := logger.NewDefaultLogger()
	log.Error(ctx, msg, zap.String("error_key", "error_value"), zap.Int("code", 500))
}

// fmt logger 實現
type fmtLogger struct{}

func (receiver fmtLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	// 構建字段字符串
	fieldStr := ""
	if ctx != nil {
		if traceID := ctx.Value(consts.TraceIDKey); traceID != nil {
			fieldStr += fmt.Sprintf(" traceID=%v", traceID)
		}
	}
	for _, f := range fields {
		fieldStr += fmt.Sprintf(" %s=%v", f.Key, f.String)
	}

	fmt.Printf("[INFO] %s%s\n", msg, fieldStr)
}

func (receiver fmtLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	// 構建字段字符串
	fieldStr := ""
	if ctx != nil {
		if traceID := ctx.Value(consts.TraceIDKey); traceID != nil {
			fieldStr += fmt.Sprintf(" traceID=%v", traceID)
		}
	}
	for _, f := range fields {
		fieldStr += fmt.Sprintf(" %s=%v", f.Key, f.String)
	}

	fmt.Printf("[ERROR] %s%s\n", msg, fieldStr)
}

func TestFmtLogger_Info(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, consts.TraceIDKey, "test_trace_id")

	logger.ChangeLogger(fmtLogger{})
	logger.Info(ctx, "test info message", zap.String("test_key", "test_value"))
}

func TestFmtLogger_Error(t *testing.T) {
	ctx := context.Background()
	msg := "test error message"

	logger.ChangeLogger(fmtLogger{})

	logger.Error(ctx, msg, zap.String("error_key", "error_value"))
}
