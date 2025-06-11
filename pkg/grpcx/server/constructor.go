package server

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"go-tool/pkg/grpcx/server/interceptor"
	"go-tool/pkg/grpcx/server/logger"
)

// Params 注入所需的依賴
type Params struct {
	fx.In

	Config   Config
	Services []IService `group:"services"`
}

func New(lc fx.Lifecycle, params Params) {
	// 創建 gRPC server 並添加中間件
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryServerRecoveryInterceptor(),
			interceptor.RequestInterceptor(),
		),
	)

	if params.Config.Logger != nil {
		logger.Init(params.Config.Logger)
	}

	// 註冊所有服務
	for _, service := range params.Services {
		service.Register(server)
	}

	// 註冊生命週期鉤子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", params.Config.Port))
			if err != nil {
				return fmt.Errorf("grpcserver.New.error: failed to listen: %v", err)
			}

			// 非阻塞啟動服務器
			go func() {
				msg := fmt.Sprintf("grpcserver.New: Server is running on port %d", params.Config.Port)
				logger.Log.Info(msg)
				if err = server.Serve(listener); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// 建立關閉超時 context
			ctx, cancel := context.WithTimeout(ctx, params.Config.ShutdownTimeout)
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
				logger.Log.Info(fmt.Sprintf("grpcserver.New: Server is shutting down, shutdownTimeout: %v", params.Config.ShutdownTimeout))
				server.Stop()
				return ctx.Err()
			case <-done:
				logger.Log.Info("grpcserver.New: Server is GracefulStop")
				return nil
			}
		},
	})
}
