package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"go-tool/pkg/logger"
	"go-tool/pkg/logger/consts"
)

func New(config Config) *slog.Logger {
	writers := []io.Writer{
		os.Stdout, // 總是輸出到控制台
	}

	if config.EnableWriteFile {
		if config.Path == "" {
			panic("slog.New config.Path is required")
		}
		// 创建 lumberjack logger
		lumberjackLogger := &lumberjack.Logger{
			Filename:  config.Path,
			MaxSize:   config.MaxSize, // megabytes
			MaxAge:    config.MaxAge,  // days
			LocalTime: true,           // 使用本地时间
			Compress:  true,           // 压缩旧文件
		}

		writers = append(writers, lumberjackLogger)
	}

	// 创建 slog handler 选项
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				timeVal := a.Value.Time()
				_, offset := timeVal.Zone()
				utcOffset := offset / 3600

				return slog.Group("",
					slog.Int64("timestamp", timeVal.UnixMilli()),
					slog.String("datetime", timeVal.Format(time.DateTime)),
					slog.String("timezone", fmt.Sprintf("UTC%+d", utcOffset)),
				)
			}

			return a
		},
	}

	l := slog.New(&CustomHandler{Handler: slog.NewJSONHandler(io.MultiWriter(writers...), opts)})
	logger.Init(l)

	return l
}

type CustomHandler struct {
	slog.Handler
}

func (receiver *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID := ctx.Value(consts.TraceIDKey); traceID != nil {
		r.Add(slog.String(consts.TraceIDKey, fmt.Sprint(traceID)))
	}
	return receiver.Handler.Handle(ctx, r)
}
