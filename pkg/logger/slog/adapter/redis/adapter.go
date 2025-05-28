package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"go-tool/pkg/logger"
)

type Logger interface {
	// Printf 給 redis.SetLogger(config.Logger);
	/*
		    可以看到這些log
			ERROR redis internal: connection lost
			WARN redis internal: connection timeout
	*/
	Printf(ctx context.Context, format string, v ...interface{})

	/*
		redis.Hook; 可以看到每個命令執行
	*/
	redis.Hook
}

// Adapter GORM logger 適配器
type Adapter struct {
	config Config
}

// New 創建 redis logger 適配器
func New(config Config) Logger {
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

func (receiver *Adapter) Printf(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	receiver.config.Logger.InfoContext(
		ctx,
		"redis.Printf",
		slog.String("message", msg),
	)
}

func (receiver *Adapter) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := next(ctx, network, addr)
		if err != nil {
			receiver.config.Logger.ErrorContext(
				ctx,
				"redis.DialHook next error",
				slog.String("network", network),
				slog.String("addr", addr),
				slog.Any("error", err),
			)
		}
		return conn, err
	}
}

func (receiver *Adapter) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmd)
		duration := time.Since(start)

		// 構建日誌字段
		fields := []any{
			slog.String("cmd", cmd.String()),
			slog.String("duration", duration.String()),
		}

		// 錯誤情況
		if err != nil {
			// redis.Nil 不是真正的錯誤
			if !errors.Is(err, redis.Nil) {
				fields = append(fields, slog.Any("error", err))
				receiver.config.Logger.ErrorContext(
					ctx,
					"redis.ProcessHook next error",
					fields...,
				)
			}
			return err
		}

		// 慢查詢
		if duration > receiver.config.SlowThreshold {
			fields = append(fields, slog.String("slow_threshold", receiver.config.SlowThreshold.String()))
			receiver.config.Logger.InfoContext(
				ctx,
				"redis.ProcessHook slow",
				fields...,
			)
			return err
		}

		// 成功的命令
		if receiver.config.EnableSuccessLog {
			receiver.config.Logger.InfoContext(
				ctx,
				"redis.ProcessHook success",
				fields...,
			)
		}

		return err
	}
}

func (receiver *Adapter) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmds)
		duration := time.Since(start)

		cmdCount := len(cmds)

		// 構建日誌字段
		fields := []any{
			slog.Int("cmdCount", cmdCount),
			slog.String("duration", duration.String()),
			slog.Any("cmds", cmds),
		}

		// 錯誤情況
		if err != nil {
			fields = append(fields, slog.Any("error", err))
			receiver.config.Logger.ErrorContext(
				ctx,
				"redis.ProcessPipelineHook error",
				fields...,
			)
			return err
		}

		// 慢管道查詢
		if duration > receiver.config.SlowThreshold {
			receiver.config.Logger.InfoContext(
				ctx,
				"redis.ProcessPipelineHook slow",
				fields...,
			)
			return err
		}

		// 成功的管道
		if receiver.config.EnableSuccessLog {
			receiver.config.Logger.InfoContext(
				ctx,
				"redis.ProcessPipelineHook success",
				fields...,
			)
		}

		return err
	}
}
