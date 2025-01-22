package goroutine_util

import (
	"context"
	"log/slog"
)

func GoWithRecover(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("[GoWithRecover]Recovered", slog.Any("error", r))
			}
		}()
		f()
	}()
}

// CopyContextKeys
// 從來源 context (src) 複製指定的 key-value 到目標 context (dst)  只複製指定的 keys，而不是全部內容
func CopyContextKeys(src, dst context.Context, keys ...string) context.Context {
	for _, key := range keys {
		if v := src.Value(key); v != nil {
			dst = context.WithValue(dst, key, v)
		}
	}

	return dst
}
