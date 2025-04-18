package ginx

import (
	"go.uber.org/fx"

	"go-tool/pkg/ginx/consts"
)

// RegisterRouter 注册路由
func RegisterRouter(constructor any) fx.Option {
	return fx.Provide(fx.Annotate(constructor, fx.As(new(IRouter)), fx.ResultTags(consts.RouterTags)))
}

// Module 注入Web 服务器
func Module(config Config) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))
	options = append(options, fx.Invoke(NewServer))

	return fx.Module("gin", options...)
}
