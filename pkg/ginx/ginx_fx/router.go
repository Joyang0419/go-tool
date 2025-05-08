package ginx_fx

import (
	"go.uber.org/fx"

	"go-tool/pkg/ginx"
	"go-tool/pkg/ginx/ginx_consts"
)

func RegisterRouter(constructor any) fx.Option {
	return fx.Provide(fx.Annotate(constructor, fx.As(new(ginx.IRouter)), fx.ResultTags(ginx_consts.RouterTags)))
}
