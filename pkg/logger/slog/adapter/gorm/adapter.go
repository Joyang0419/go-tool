package gorm

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	gormLogger "gorm.io/gorm/logger"

	"go-tool/pkg/logger"
)

// Adapter GORM logger 適配器
type Adapter struct {
	config Config
}

// New 創建 GORM logger 適配器
func New(config Config) gormLogger.Interface {
	adapter := new(Adapter)
	if config.Logger == nil {
		config.Logger = logger.Log
	}
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 500 * time.Millisecond
	}

	adapter.config = config

	return adapter
}

func (receiver *Adapter) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	// 創建新的適配器實例
	return receiver
}

func (receiver *Adapter) Info(ctx context.Context, msg string, args ...interface{}) {
	receiver.config.Logger.InfoContext(ctx, msg, "gorm_info", args)
}

func (receiver *Adapter) Warn(ctx context.Context, msg string, args ...interface{}) {
	receiver.config.Logger.InfoContext(ctx, msg, "gorm_warn", args)
}

func (receiver *Adapter) Error(ctx context.Context, msg string, args ...interface{}) {
	receiver.config.Logger.ErrorContext(ctx, msg, "gorm_error", args)
}

func (receiver *Adapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	// 構建日誌字段
	fields := []any{
		slog.String("elapsed", elapsed.String()),
		slog.String("sql", sql),
		slog.Int64("rowsAffected", rowsAffected),
	}

	switch {
	case err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound):
		// 有錯誤且不是 RecordNotFound
		receiver.config.Logger.ErrorContext(ctx, "GORM SQL Error",
			append(fields, "error", err.Error())...)

	case elapsed > receiver.config.SlowThreshold:
		// 慢查詢：超過閾值
		receiver.config.Logger.InfoContext(ctx, "GORM Slow SQL",
			append(fields, "slow_threshold", receiver.config.SlowThreshold.String())...)

	default:
		if receiver.config.EnableSuccessSQLLog {
			receiver.config.Logger.InfoContext(ctx, "GORM SQL", fields...)
		}
	}
}
