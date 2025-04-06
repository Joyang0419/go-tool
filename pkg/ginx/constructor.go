package ginx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"golang.org/x/exp/slices"

	"go-tool/pkg/ginx/middleware"
)

// Params 注入所需的依賴
type Params struct {
	fx.In

	Config  Config
	Routers []IRouter `group:"routers"`
}

func NewServer(lc fx.Lifecycle, params Params) {
	allowModes := []string{gin.DebugMode, gin.ReleaseMode, gin.TestMode}
	if !slices.Contains(allowModes, params.Config.Mode) {
		panic(fmt.Sprintf("ginx.NewServer gin mode %s is not allow, only allow %v", params.Config.Mode, allowModes))
	}

	gin.SetMode(params.Config.Mode)

	engine := gin.Default()

	// recovery 必須放在最前面
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.TraceIDMiddleware())

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
