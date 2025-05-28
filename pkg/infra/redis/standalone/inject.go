package standalone

import (
	"go-tool/pkg/infra/redis/logger"

	"go.uber.org/fx"
)

func FXModule(config Config, l ...logger.Interface) fx.Option {
	var options []fx.Option
	if len(l) > 0 {
		config.Logger = l[0]
	}

	options = append(options, fx.Supply(config))
	options = append(options, fx.Provide(New))

	return fx.Module("redis_standalone_module", options...)
}
