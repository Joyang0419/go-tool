package inject

import (
	"go.uber.org/fx"

	"go-tool/pkg/grpcx/consts"
	"go-tool/pkg/grpcx/server"
)

// RegisterService 注册服务
func RegisterService(constructor any) fx.Option {
	return fx.Provide(fx.Annotate(constructor, fx.As(new(server.IService)), fx.ResultTags(consts.ServiceTags)))
}

type FXModuleParam struct {
	ServiceConstructors []any
}

func FXModule(config server.Config, param FXModuleParam) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))

	for _, constructor := range param.ServiceConstructors {
		options = append(options, RegisterService(constructor))
	}

	options = append(options, fx.Invoke(server.New))

	return fx.Module("grpcserver_module", options...)
}
