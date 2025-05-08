package ginx_fx

import (
	"go.uber.org/fx"

	"go-tool/pkg/ginx"
	"go-tool/pkg/ginx/ginx_consts"
	"go-tool/pkg/ginx/ginx_middleware"
)

func InjectCORSMiddleware(config ginx_middleware.CORSConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Provide(fx.Annotate(ginx_middleware.NewCORSMiddleware, fx.As(new(ginx.IMiddleware)), fx.ResultTags(ginx_consts.MiddlewareTags))),
	)
}
