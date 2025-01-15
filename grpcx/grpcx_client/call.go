package grpcx_client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"

	pkgerrors "github.com/pkg/errors"

	"go-tool/ginx/consts"
)

/*
callOpt 應用
WaitForReady = false (默認)

如果服務器暫時不可用，立即返回錯誤
錯誤類型通常是 UNAVAILABLE
客戶端需要自己處理重試邏輯


WaitForReady = true

如果服務器暫時不可用，會等待直到服務恢復
在底層實現了重試機制
適合對延遲不敏感但需要確保請求成功的場景
*/

func CallByClient[C any, Req any, Resp any](
	ctx context.Context,
	conn *grpc.ClientConn,
	newClient func(grpc.ClientConnInterface) C,
	method func(C, context.Context, Req, ...grpc.CallOption) (Resp, error),
	request Req,
	callOpts ...grpc.CallOption,
) (response Resp, err error) {
	if conn == nil {
		return response, pkgerrors.New("[CallByClient]conn is nil")
	}
	// 檢查連接狀態
	state := conn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		return response, fmt.Errorf("[CallByClient]connection is not ready, current state: %s", state)
	}

	// 設置 traceID
	ctx = SetupTraceID(ctx)
	client := newClient(conn)

	return method(client, ctx, request, callOpts...)
}

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
