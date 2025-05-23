package interceptor

import (
	"context"
	"log/slog"
	"runtime/debug"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				// 從堆疊中找出最相關的位置
				var location string
				// 獲取堆疊跟踪
				stack := debug.Stack()

				// 解析堆疊資訊找到程式碼位置
				lines := strings.Split(string(stack), "\n")

				for _, line := range lines {
					if strings.Contains(line, ".go:") &&
						!strings.Contains(line, "runtime/") &&
						!strings.Contains(line, "grpc/") &&
						!strings.Contains(line, ".pb.go") &&
						!strings.Contains(line, "go-go-tool") {
						if idx := strings.Index(line, " +0x"); idx != -1 {
							location = strings.TrimSpace(line[:idx])
						} else {
							location = strings.TrimSpace(line)
						}
						break
					}
				}

				// 記錄日誌
				slog.ErrorContext(ctx, "[UnaryServerRecoveryInterceptor]Panic recovered",
					slog.Any("error", r),
					slog.String("location", location),
					slog.String("method", info.FullMethod),
				)

				err = status.Errorf(codes.Internal, "Internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
