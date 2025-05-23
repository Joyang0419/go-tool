package fx

import (
	"go.uber.org/fx"

	"go-tool/pkg/httpclient"
	"go-tool/pkg/httpclient/resty"
)

// Module 注入Web 服务器
func Module(config httpclient.Config) fx.Option {
	var options []fx.Option
	options = append(options, fx.Supply(config))
	options = append(options, fx.Provide(resty.NewHTTPClient))

	return fx.Module("httpclient_module", options...)
}
