package client

import (
	"context"

	"google.golang.org/grpc/metadata"

	"go-tool/pkg/grpcx/consts"
)

// SetupTraceID 將 traceID 加入到 outgoing context
func SetupTraceID(ctx context.Context) context.Context {
	// 從 context 中獲取 traceID
	traceID, ok := ctx.Value(consts.TraceIDKey).(string)
	if !ok || traceID == "" {
		return ctx
	}

	// 檢查是否已經存在 metadata
	var md metadata.MD
	if existingMD, exists := metadata.FromOutgoingContext(ctx); exists {
		// 如果已有 metadata，複製並添加 traceID
		md = existingMD.Copy()
	} else {
		// 如果沒有現有的 metadata，創建新的
		md = metadata.New(nil)
	}

	// 設置 traceID
	md.Set(consts.TraceIDKey, traceID)

	// 創建新的 outgoing context
	return metadata.NewOutgoingContext(ctx, md)
}
