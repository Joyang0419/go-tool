package grpcx_server

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"go-tool/grpcx/grpcx_server/binding/logger"

	"go-tool/grpcx/grpcx_server/interceptor"
)

type ServerConfig struct {
	Port            int
	ShutdownTimeout time.Duration
	Logger          logger.ILogger
}

// ServerParams 注入所需的依賴
type ServerParams struct {
	fx.In

	ServerConfig ServerConfig
	Services     []IService `group:"services"`
}

func NewServer(lc fx.Lifecycle, params ServerParams) {
	// 創建 gRPC server 並添加中間件
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryServerRecoveryInterceptor(),
			interceptor.RequestInterceptor(),
		),
	)

	// 註冊所有服務
	for _, service := range params.Services {
		service.Register(server)
	}

	// 設置 logger
	if params.ServerConfig.Logger != nil {
		logger.ChangeLogger(params.ServerConfig.Logger)
	}

	// 註冊生命週期鉤子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", params.ServerConfig.Port))
			if err != nil {
				return fmt.Errorf("failed to listen: %v", err)
			}

			// 非阻塞啟動服務器
			go func() {
				logger.Info(context.TODO(), fmt.Sprintf("[NewServer]Server is running on port %d", params.ServerConfig.Port))
				if err = server.Serve(listener); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// 建立關閉超時 context
			ctx, cancel := context.WithTimeout(ctx, params.ServerConfig.ShutdownTimeout)
			defer cancel()

			// 創建一個 channel 來接收關閉完成信號
			done := make(chan struct{})
			go func() {
				server.GracefulStop()
				close(done)
			}()

			// 等待關閉完成或超時
			select {
			case <-ctx.Done():
				logger.Info(context.TODO(), fmt.Sprintf("[NewServer]Server is shutting down, shutdownTimeout: %v", params.ServerConfig.ShutdownTimeout))
				server.Stop()
				return ctx.Err()
			case <-done:
				logger.Info(context.TODO(), "[NewServer]Server is GracefulStop")
				return nil
			}
		},
	})
}
