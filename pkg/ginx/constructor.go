package ginx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"golang.org/x/exp/slices"
)

// Params 注入所需的依賴
type Params struct {
	fx.In

	Config      Config
	Routers     []IRouter     `group:"routers"`
	Middlewares []IMiddleware `group:"middlewares"`
}

func NewServer(lc fx.Lifecycle, params Params) {
	allowModes := []string{gin.DebugMode, gin.ReleaseMode, gin.TestMode}
	if !slices.Contains(allowModes, params.Config.Mode) {
		panic(fmt.Sprintf("ginx.NewServer gin mode is not allow, only allow %v, mode: %v", allowModes, params.Config.Mode))
	}

	gin.SetMode(params.Config.Mode)

	engine := gin.Default()

	for _, middleware := range params.Middlewares {
		engine.Use(middleware.HandlerFunc())
	}

	for _, router := range params.Routers {
		router.Routes(engine)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", params.Config.Port),
		Handler: engine,
	}

	// 註冊生命週期鉤子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 非阻塞啟動服務器
			go func() {
				slog.Info(fmt.Sprintf("ginx.NewServer Server is running on port %d", params.Config.Port))
				if err := server.ListenAndServe(); err != nil {
					slog.Error(fmt.Sprintf("ginx.NewServer server.ListenAndServe() error: %v", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info(fmt.Sprintf("ginx.NewServer Server is shutting down, shutdownTimeout: %v", params.Config.ShutdownTimeout))
			ctx, cancel := context.WithTimeout(ctx, params.Config.ShutdownTimeout)
			defer cancel()

			return server.Shutdown(ctx)
		},
	})
}
