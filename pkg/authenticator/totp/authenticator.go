package totp

/*
	Time-based One-Time Password
*/

import (
	"github.com/pquerna/otp/totp"

	"go-tool/pkg/authenticator"
)

type Authenticator struct{}

type SO struct {
	SecretKey string
	InputCode string
}

func New() *Authenticator {
	return &Authenticator{}
}

func (receiver *Authenticator) Verify(so SO) error {
	valid := totp.Validate(so.InputCode, so.SecretKey)
	if !valid {
		return authenticator.ErrVerifyFailed
	}

	return nil
}
