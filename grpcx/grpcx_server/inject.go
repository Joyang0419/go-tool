package grpcx_server

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"go-tool/grpcx/consts"
)

type IService interface {
	Register(*grpc.Server)
}

// RegisterService 注册服务
func RegisterService(constructor any) fx.Option {
	return fx.Provide(fx.Annotate(constructor, fx.As(new(IService)), fx.ResultTags(consts.ServiceTags)))
}

// InjectServer 注入Web 服务器
func InjectServer(config ServerConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Invoke(NewServer),
	)
}
