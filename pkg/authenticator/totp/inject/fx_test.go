package inject

import (
	"fmt"
	"testing"

	"go.uber.org/fx"

	"go-tool/pkg/authenticator/totp"
)

func TestFXModule(t *testing.T) {
	fx.New(
		FXModule(),
		fx.Invoke(
			func(authenticator *totp.Authenticator) {
				err := authenticator.Verify(totp.SO{
					SecretKey: "T6YTKACIDCFN7DGLKDT76FJ52KYTCBC3",
					InputCode: "992414",
				})
				fmt.Println(err)
			},
		))
}
