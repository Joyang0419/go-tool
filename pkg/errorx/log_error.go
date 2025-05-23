package errorx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

func LogError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	stackTrace := StackTrace(err)
	stackTraceLen := len(stackTrace)
	if stackTraceLen == 0 {
		slog.ErrorContext(
			ctx,
			"LogError stackTraceLen is zero",
			slog.Any("error", err),
			slog.Any("stackTrace", stackTrace),
		)
		return
	}
	lastStackTrace := stackTrace[0]
	// full info; pkgerrors formater
	/*
		%s - 檔案名稱 (base name)
		%+s - 完整函數名稱和檔案路徑
		%d - 行號
		%n - 函數名稱
		%v - 相當於 %s:%d (檔案:行號)
		%+v - 相當於 %+s:%d (完整資訊)
	*/
	msg := fmt.Sprintf("%+v.error", lastStackTrace)
	// 方法1: 建立新的 logger，不使用 CallerSkip
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	// 直接建立原生 handler，不走 slogx
	handler := slog.NewJSONHandler(os.Stdout, opts)
	_ = slog.New(handler)

	// 手動設定正確的 caller
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // 這次應該正確了
	record := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
	record.Add(slog.Any("error", err))

	// 從 context 取得 traceID (如果你的 slogx 有這個功能)
	if traceID := ctx.Value("traceID"); traceID != nil {
		record.Add(slog.String("traceID", traceID.(string)))
	}

	_ = handler.Handle(ctx, record)
}
