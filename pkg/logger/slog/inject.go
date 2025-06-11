package slog

import (
	"log/slog"

	"go.uber.org/fx"
)

// FXModule 注入Slog
func FXModule(config Config) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))
	options = append(options, fx.Provide(New))
	options = append(options, fx.Invoke(slog.SetDefault))

	return fx.Module("slog_module", options...)
}
