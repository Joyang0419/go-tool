package concurrency

import (
	"log/slog"

	"github.com/panjf2000/ants/v2"
	"github.com/sourcegraph/conc/iter"
)

// GoWithRecovery 並帶有錯誤恢復機制GoRoutine，避免Panic導致服務關閉
func GoWithRecovery(fn func()) {
	var opts []ants.Option
	opts = append(opts, ants.WithPanicHandler(func(i interface{}) {
		slog.Error("GoWithRecovery found panic", slog.Any("panic", i))
	}))

	p, _ := ants.NewPool(1, opts...)

	if err := p.Submit(fn); err != nil {
		slog.Error("GoWithRecovery submit error", slog.Any("error", err))
	}
}

// GoForSameTasks 「對輸入陣列中的每個元素應用指定的函式，並返回一個包含結果的新陣列。」
func GoForSameTasks[TInput any, TResult any](inputs []TInput, fn func(input *TInput) (TResult, error)) ([]TResult, error) {
	return iter.MapErr(inputs, fn)
}
