package ginx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"go-tool/pkg/ginx/middleware"
)

type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// ServerParams 注入所需的依賴
type ServerParams struct {
	fx.In

	ServerConfig ServerConfig
	Routers      []IRouter `group:"routers"`
}

func NewServer(lc fx.Lifecycle, params ServerParams) {
	engine := gin.Default()

	// recovery 必須放在最前面
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.TraceIDMiddleware())

	for _, router := range params.Routers {
		router.Routes(engine)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", params.ServerConfig.Port),
		Handler: engine,
	}

	// 註冊生命週期鉤子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 非阻塞啟動服務器
			go func() {
				slog.Info(fmt.Sprintf("[NewServer]Server is running on port %d", params.ServerConfig.Port))
				if err := server.ListenAndServe(); err != nil {
					slog.Error(fmt.Sprintf("[NewServer]server.ListenAndServe() error: %v", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info(fmt.Sprintf("[NewServer]Server is shutting down, shutdownTimeout: %v", params.ServerConfig.ShutdownTimeout))
			ctx, cancel := context.WithTimeout(ctx, params.ServerConfig.ShutdownTimeout)
			defer cancel()

			return server.Shutdown(ctx)
		},
	})
}
