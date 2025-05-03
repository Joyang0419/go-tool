package setup

import (
	"net/http"

	"go-tool/pkg/ginx/ginx_error"
	"go-tool/pkg/ginx/ginx_error/default_error"
)

var (
	ErrPanic    ginx_error.IGinxError = default_error.DefaultError{StatusCode: http.StatusInternalServerError, CustomCode: default_error.CustomCodeServerSideServerError}
	ErrParsedSO ginx_error.IGinxError = default_error.DefaultError{StatusCode: http.StatusBadRequest, CustomCode: default_error.CustomCodeClientSideBadRequest}
)
