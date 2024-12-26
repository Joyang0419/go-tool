package web

import (
	"go.uber.org/fx"
)

// RegisterController 注册控制器
func RegisterController(constructor any) fx.Option {
	return fx.Provide(fx.Annotate(constructor, fx.As(new(IController)), fx.ResultTags(`group:"controllers"`)))
}

// InjectServer 注入Web 服务器
func InjectServer(config ServerConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Invoke(NewServer),
	)
}
