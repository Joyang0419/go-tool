package ginx

import (
	"go.uber.org/fx"

	"go-tool/pkg/ginx/consts"
)

// RegisterRouter 注册路由
func RegisterRouter(constructor any) fx.Option {
	return fx.Provide(fx.Annotate(constructor, fx.As(new(IRouter)), fx.ResultTags(consts.RouterTags)))
}

// InjectServer 注入Web 服务器
func InjectServer(config ServerConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Invoke(NewServer),
	)
}
