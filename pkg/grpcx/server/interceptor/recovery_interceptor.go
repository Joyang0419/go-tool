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
				// 獲取實際發生 panic 的位置
				location := getPanicLocation()

				// 記錄日誌
				slog.ErrorContext(ctx, "interceptor.UnaryServerRecoveryInterceptor.error: panic recovered",
					slog.Any("error", r),
					slog.String("location", location),
					slog.String("method", info.FullMethod),
					slog.String("stack", string(debug.Stack())), // 完整堆疊供調試用
				)

				err = status.Errorf(codes.Internal, "Internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// getPanicLocation 從堆疊中提取實際發生 panic 的位置
func getPanicLocation() string {
	stack := debug.Stack()
	lines := strings.Split(string(stack), "\n")

	// Go 堆疊格式：
	// goroutine 1 [running]:
	// function_name(args...)
	//     /path/to/file.go:line +0xoffset
	// next_function(args...)
	//     /path/to/file.go:line +0xoffset

	for i := 0; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		nextLine := strings.TrimSpace(lines[i+1])

		// 檢查下一行是否是檔案位置（包含 .go: 且不是框架代碼）
		if strings.Contains(nextLine, ".go:") &&
			!strings.Contains(nextLine, "runtime/") &&
			!strings.Contains(nextLine, "grpc/") &&
			!strings.Contains(nextLine, ".pb.go") &&
			!strings.Contains(nextLine, "go-tool/pkg/grpcx") &&
			!strings.Contains(nextLine, "go-go-tool") &&
			!strings.Contains(line, "panic(") &&
			!strings.Contains(line, "recovery") {

			// 提取檔案位置（移除 +0x 部分）
			location := nextLine
			if idx := strings.Index(location, " +0x"); idx != -1 {
				location = location[:idx]
			}

			return strings.TrimSpace(location)
		}
	}

	return "unknown location"
}
