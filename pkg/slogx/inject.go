package slogx

import (
	"log/slog"

	"go.uber.org/fx"
)

// Module 注入Slog
func Module(config Config) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))
	options = append(options, fx.Provide(NewSlog))
	options = append(options, fx.Invoke(slog.SetDefault))

	return fx.Module("slog", options...)
}
