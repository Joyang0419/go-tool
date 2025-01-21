package slogx

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewSlog(config Config) *slog.Logger {
	writers := []io.Writer{
		os.Stdout, // 總是輸出到控制台
	}

	// 如果有path, 則輸出到檔案
	if config.Path != "" {
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
			case slog.SourceKey:
				if source, ok := a.Value.Any().(*slog.Source); ok {
					callers := make([]uintptr, 32)
					n := runtime.Callers(0, callers)
					frames := runtime.CallersFrames(callers[:n])

					// 跳過指定的幀數
					for i := 0; i <= config.CallerSkip; i++ {
						_, more := frames.Next()
						if !more {
							break
						}
					}

					// 獲取正確的調用幀
					if frame, more := frames.Next(); more {
						source.File = frame.File
						source.Line = frame.Line
						source.Function = frame.Function
					}
				}

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

	return slog.New(&CustomHandler{Handler: slog.NewJSONHandler(io.MultiWriter(writers...), opts)})
}

type CustomHandler struct {
	slog.Handler
}

func (receiver *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		r.Add(slog.String(TraceIDKey, fmt.Sprint(traceID)))
	}
	return receiver.Handler.Handle(ctx, r)
}
