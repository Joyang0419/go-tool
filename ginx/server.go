package ginx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"go-tool/ginx/binding/logger"
	"go-tool/ginx/middleware"
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
	Routers      []IRouter `group:"routers"`
}

func NewServer(lc fx.Lifecycle, params ServerParams) {
	engine := gin.Default()
	engine.Use(
		middleware.RecoveryMiddleware(),
		middleware.TraceIDMiddleware(),
	)

	for _, router := range params.Routers {
		router.Routes(engine)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", params.ServerConfig.Port),
		Handler: engine,
	}
	if params.ServerConfig.Logger != nil {
		logger.ChangeLogger(params.ServerConfig.Logger)
	}

	// 註冊生命週期鉤子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 非阻塞啟動服務器
			go func() {
				logger.Info(context.TODO(), fmt.Sprintf("[NewServer]Server is running on port %d", params.ServerConfig.Port))
				if err := server.ListenAndServe(); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info(context.TODO(), fmt.Sprintf("[NewServer]Server is shutting down, shutdownTimeout: %v", params.ServerConfig.ShutdownTimeout))
			ctx, cancel := context.WithTimeout(ctx, params.ServerConfig.ShutdownTimeout)
			defer cancel()

			return server.Shutdown(ctx)
		},
	})
}
