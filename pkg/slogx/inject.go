package slogx

import (
	"log/slog"

	"go.uber.org/fx"
)

// InjectSlog 注入Slog
func InjectSlog(config Config) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Provide(NewSlog),
		fx.Invoke(slog.SetDefault),
	)
}
