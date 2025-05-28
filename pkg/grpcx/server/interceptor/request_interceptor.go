package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go-tool/pkg/grpcx/consts"
)

// RequestInterceptor 返回一個一元 RPC 攔截器
func RequestInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()
		trace := traceID(ctx)

		if trace != "" {
			ctx = context.WithValue(ctx, consts.TraceIDKey, trace)
		}

		// 準備記錄請求
		slog.InfoContext(ctx, "[RequestInterceptor]Grpc incoming request",
			slog.String("method", path.Base(info.FullMethod)),
			slog.Any("request", req),
			slog.String("traceID", trace),
		)

		// 執行 handler
		resp, err := handler(ctx, req)

		// 計算處理時間
		latency := time.Since(startTime)

		// 準備日誌字段
		var fields []any
		fields = append(fields, slog.String("method", path.Base(info.FullMethod)))
		fields = append(fields, slog.Any("latencySeconds", latency))
		fields = append(fields, slog.String("traceID", trace))

		// 如果有錯誤，添加錯誤資訊
		if err != nil {
			fields = append(fields, slog.Any("request", req))
			fields = append(fields, slog.Any("response", resp))

			st, ok := status.FromError(err)
			if ok {
				fields = append(fields, slog.String("errorCode", st.Code().String()))
				fields = append(fields, slog.String("errorMessage", st.Message()))

				slog.ErrorContext(ctx, fmt.Sprintf("[RequestInterceptor]Grpc request failed: error: %v", err), fields...)
				return resp, err
			}

			slog.ErrorContext(ctx, fmt.Sprintf("[RequestInterceptor]Grpc request failed error: %v", err), fields...)
			return resp, err
		}

		// 記錄成功請求
		slog.InfoContext(ctx, "[RequestInterceptor]Grpc request success", fields...)

		return resp, nil
	}
}

// traceID 從 metadata 中提取 traceID
func traceID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(consts.TraceIDKey)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}
