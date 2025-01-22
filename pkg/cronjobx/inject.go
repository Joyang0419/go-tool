package cronjobx

import (
	"go.uber.org/fx"
)

// RegisterJob 註冊工作
func RegisterJob(constructor any) fx.Option {
	return fx.Provide(fx.Annotate(constructor, fx.As(new(IJob)), fx.ResultTags(`group:"jobs"`)))
}

// InjectCronJob 注入Web 服务器
func InjectCronJob(config Config) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Invoke(NewManager),
	)
}
