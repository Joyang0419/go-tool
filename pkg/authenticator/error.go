package authenticator

import (
	pkgerrors "github.com/pkg/errors"
)

var (
	ErrVerifyFailed = pkgerrors.New("authenticator verify failed")
)
