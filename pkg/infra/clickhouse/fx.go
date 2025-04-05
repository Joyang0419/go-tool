package clickhouse

import (
	"go.uber.org/fx"
)

func Module(config Config) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))
	options = append(options, fx.Provide(NewConnection))

	return fx.Module("clickhouse", options...)
}
