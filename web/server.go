package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"go-tool/web/middleware"
)

type ServerConfig struct {
	Port            int
	ShutdownTimeout time.Duration
}

func NewServerConfig(port int, shutdownTimeout time.Duration) ServerConfig {
	return ServerConfig{
		Port:            port,
		ShutdownTimeout: shutdownTimeout,
	}
}

// ServerParams 注入所需的依賴
type ServerParams struct {
	fx.In

	ServerConfig ServerConfig
	Controllers  []TIController `group:"controllers"`
}

func NewServer(lc fx.Lifecycle, params ServerParams) {
	engine := gin.Default()
	engine.Use(middleware.Recovery())

	for _, controller := range params.Controllers {
		group := engine.Group(controller.BasePath())
		group.Use(controller.Middlewares()...)
		controller.Routes(group)
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
				if err := server.ListenAndServe(); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, params.ServerConfig.ShutdownTimeout)
			defer cancel()

			return server.Shutdown(ctx)
		},
	})
}
