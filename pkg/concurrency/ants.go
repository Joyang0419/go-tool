package concurrency

import (
	"context"
	"log/slog"

	"github.com/panjf2000/ants/v2"
	"github.com/samber/lo"
)

// WithRecovery 並帶有錯誤恢復機制GoRoutine，避免Panic導致服務關閉
func WithRecovery(ctx context.Context, fn func()) {
	panicHandler := func(r interface{}) {
		slog.ErrorContext(ctx, "WithRecovery panic recovered", slog.Any("panic", r))
	}

	p, err := ants.NewPool(0, ants.WithPanicHandler(panicHandler))
	if !lo.IsNil(err) {
		slog.ErrorContext(ctx, "WithRecovery ants.NewPool error", slog.Any("error", err.Error()))
		return
	}
	defer p.Release()

	if err = p.Submit(fn); !lo.IsNil(err) {
		slog.ErrorContext(ctx, "WithRecovery p.Submit error", slog.Any("error", err.Error()))
		return
	}
}
