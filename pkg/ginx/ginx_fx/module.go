package ginx_fx

import (
	"go.uber.org/fx"

	"go-tool/pkg/ginx"
)

// Module 注入Web 服务器
func Module(config ginx.Config) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))
	options = append(options, fx.Invoke(ginx.NewServer))

	return fx.Module("ginx_module", options...)
}
