package grpcx_client

import (
	"context"
	"fmt"
	"time"

	"github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"framework/pkg/ginx/consts"
)

func Call[C any, Req any, Resp any](
	ctx context.Context,
	pool *grpcpool.Pool,
	newClient func(grpc.ClientConnInterface) C, // 改用 ClientConnInterface
	method func(C, context.Context, Req, ...grpc.CallOption) (Resp, error), // 加入 CallOption
	request Req,
	timeout ...time.Duration,
) (response Resp, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[grpcx][Call]panic error: %v", r)
		}
	}()
	conn, err := pool.Get(ctx)
	if err != nil {
		return response, err
	}
	defer func() {
		_ = conn.Close()
	}()

	ctx = SetupTraceID(ctx)

	to := 5 * time.Second
	if len(timeout) > 0 {
		to = timeout[0]
	}

	ctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	client := newClient(conn.ClientConn)
	return method(client, ctx, request) // CallOption 是可選的，所以這裡不傳也可以
}

// SetupTraceID 將 traceID 加入到 outgoing context
func SetupTraceID(ctx context.Context) context.Context {
	// 從 context 中獲取 traceID
	traceID, ok := ctx.Value("traceId").(string)
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
