package fx

import (
	"go.uber.org/fx"

	"go-tool/pkg/ginx"
	"go-tool/pkg/ginx/consts"
	"go-tool/pkg/ginx/middleware"
)

func InjectRecoveryMiddleware() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(middleware.NewRecoveryMiddleware, fx.As(new(ginx.IMiddleware)), fx.ResultTags(consts.MiddlewareTags))),
	)
}

func InjectAccessLogMiddleware() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(middleware.NewAccessLogMiddleware, fx.As(new(ginx.IMiddleware)), fx.ResultTags(consts.MiddlewareTags))),
	)
}

func InjectErrorMiddleware() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(middleware.NewErrorMiddleware, fx.As(new(ginx.IMiddleware)), fx.ResultTags(consts.MiddlewareTags))),
	)
}

func InjectTraceIDMiddleware() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(middleware.NewTraceIDMiddleware, fx.As(new(ginx.IMiddleware)), fx.ResultTags(consts.MiddlewareTags))),
	)
}

func InjectCORSMiddleware(config middleware.CORSConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Provide(fx.Annotate(middleware.NewCORSMiddleware, fx.As(new(ginx.IMiddleware)), fx.ResultTags(consts.MiddlewareTags))),
	)
}
