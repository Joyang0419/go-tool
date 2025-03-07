package slogx

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewSlog(config Config) *slog.Logger {
	writers := []io.Writer{
		os.Stdout, // 總是輸出到控制台
	}

	if config.EnableWriteFile {
		if config.Path == "" {
			panic("NewSlog config.Path is required")
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

					frameCopy := runtime.CallersFrames(callers[:n])
					// 跳過自己的幀數
					for i := 0; i <= config.CallerSkip+1; i++ {
						_, more := frameCopy.Next()
						if !more {
							break
						}
					}

					// 獲取上層的兩個調用幀
					stackNum := 2
					stackFrames := make([]string, 0, stackNum)
					for i := 0; i < stackNum; i++ {
						if frame, more := frameCopy.Next(); more {
							stackFrames = append(stackFrames, fmt.Sprintf("%s:%d", frame.File, frame.Line))
						}
					}

					return slog.Group("",
						slog.String("file", fmt.Sprintf("%s:%d", source.File, source.Line)),
						slog.String("function", source.Function),
						slog.String("stack", strings.Join(stackFrames, ", ")),
					)
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
