package slogx

import (
	"go.uber.org/fx"

	"go-tool/pkg/logger"
)

// FXModule 注入Slog
func FXModule(config Config) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))
	options = append(options, fx.Provide(NewSlog))
	options = append(options, fx.Invoke(logger.Init))

	return fx.Module("slog", options...)
}
