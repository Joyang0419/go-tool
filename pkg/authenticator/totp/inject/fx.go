package inject

import (
	"go.uber.org/fx"

	"go-tool/pkg/authenticator/totp"
)

func FXModule() fx.Option {
	var options []fx.Option

	options = append(options, fx.Provide(totp.New))

	return fx.Module("authenticator_totp_module", options...)
}
