package logger

import (
	"context"

	"go.uber.org/zap"

	"go-tool/ginx/consts"
)

/*
	Logger
	Why抽出來? 因為直接使用其他模組的log, 如果模組的log呼叫方式改變，會很多地方要改。
	參考的想法:
		gorm.Config Logger 抽出來的想法
		gin: binding.Validator; 沒有自定義設定時，就用DefaultValidator
*/

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Error(ctx, msg, fields...)
}

type ILogger interface {
	Info(context.Context, string, ...zap.Field)
	Error(context.Context, string, ...zap.Field)
}

var logger = NewDefaultLogger()

type defaultLogger struct {
	log *zap.Logger
}

func NewDefaultLogger() ILogger {
	config := zap.NewProductionConfig()
	l, _ := config.Build(zap.AddCallerSkip(1))
	return &defaultLogger{
		log: l,
	}
}

func (receiver *defaultLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	zapFields := make([]zap.Field, 0, len(fields)+1)

	if ctx != nil {
		if traceID := ctx.Value(consts.TraceIDKey); traceID != nil {
			zapFields = append(zapFields, zap.Any(consts.TraceIDKey, traceID))
		}
	}

	zapFields = append(zapFields, fields...)
	receiver.log.Info(msg, zapFields...)
}

func (receiver *defaultLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	zapFields := make([]zap.Field, 0, len(fields)+1)

	if ctx != nil {
		if traceID := ctx.Value(consts.TraceIDKey); traceID != nil {
			zapFields = append(zapFields, zap.Any(consts.TraceIDKey, traceID))
		}
	}

	zapFields = append(zapFields, fields...)
	receiver.log.Error(msg, zapFields...)
}

func ChangeLogger(l any) ILogger {
	if customLogger, ok := l.(ILogger); ok {
		logger = customLogger
		return customLogger
	}
	panic("logger must implement ILogger interface")
}
