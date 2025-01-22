package goroutine_util

import (
	"context"
	"testing"
)

func TestCopyContextKeys(t *testing.T) {
	// 建立來源 context 並設置一些值
	srcCtx := context.Background()
	srcCtx = context.WithValue(srcCtx, "traceId", "trace-123")
	srcCtx = context.WithValue(srcCtx, "userId", "user-456")
	srcCtx = context.WithValue(srcCtx, "requestId", "req-789")

	// 建立目標 context
	dstCtx := context.Background()

	// 只複製 traceId 和 requestId
	dstCtx = CopyContextKeys(srcCtx, dstCtx, "traceId", "requestId")

	// 驗證複製的值
	if v := dstCtx.Value("traceId"); v != "trace-123" {
		t.Errorf("expected traceId=trace-123, got %v", v)
	}

	if v := dstCtx.Value("requestId"); v != "req-789" {
		t.Errorf("expected requestId=req-789, got %v", v)
	}

	// 驗證未複製的值
	if v := dstCtx.Value("userId"); v != nil {
		t.Errorf("expected userId=nil, got %v", v)
	}
}
