package ctx_util

import (
	"context"

	"github.com/google/uuid"
)

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

func WithTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, "traceId", uuid.New().String())
}
