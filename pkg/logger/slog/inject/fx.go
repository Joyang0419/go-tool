package inject

import (
	"go.uber.org/fx"

	"go-tool/pkg/logger"
	"go-tool/pkg/logger/slog"
)

// FXModule 注入Slog
func FXModule(config slogx.Config) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))
	options = append(options, fx.Provide(slogx.NewSlog))
	options = append(options, fx.Invoke(logger.Init))

	return fx.Module("slog", options...)
}
