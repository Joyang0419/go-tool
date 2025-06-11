package gorm

import (
	"context"
	"fmt"
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
		config.SlowThreshold = 5000 * time.Millisecond
	}

	adapter.config = config

	return adapter
}

func (receiver *Adapter) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	return receiver
}

func (receiver *Adapter) Info(ctx context.Context, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	finalMsg := fmt.Sprintf("gorm.Info msg: %v", formattedMsg)
	receiver.config.Logger.InfoContext(ctx, finalMsg)
}

func (receiver *Adapter) Warn(ctx context.Context, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	finalMsg := fmt.Sprintf("gorm.Warn msg: %v", formattedMsg)
	receiver.config.Logger.InfoContext(ctx, finalMsg)
}

func (receiver *Adapter) Error(ctx context.Context, msg string, args ...interface{}) {
	// 先格式化完整的訊息
	formattedMsg := fmt.Sprintf(msg, args...)
	finalMsg := fmt.Sprintf("gorm.Error msg: %v", formattedMsg)

	// 只傳訊息，不傳額外參數
	receiver.config.Logger.ErrorContext(ctx, finalMsg)
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

	// 錯誤SQL, 且是真的SQL語句錯誤，不是SQL語句正確，且沒資料的錯誤
	if err != nil {
		if !errors.Is(err, gormLogger.ErrRecordNotFound) {
			fields = append(fields, slog.Any("error", err))
			slog.ErrorContext(
				ctx,
				"gorm.Trace error",
				fields...,
			)
		}
		return
	}

	// 慢查詢：超過閾值
	if elapsed > receiver.config.SlowThreshold {
		fields = append(
			fields,
			slog.String("slow_threshold", receiver.config.SlowThreshold.String()),
		)
		slog.InfoContext(
			ctx,
			"gorm.Trace slowSQL",
			fields...,
		)
	}

	if receiver.config.EnableSuccessLog {
		receiver.config.Logger.InfoContext(
			ctx,
			"gorm.Trace success",
			fields...,
		)
	}
}
