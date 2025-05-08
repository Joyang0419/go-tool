package ginx_error

import (
	"net/http"
)

const (
	ErrCodeServerSideServerError = 500000
)

const (
	ErrCodeClientSideBadRequest = 400000
)

var (
	ErrBadParam      = DefaultError{statusCode: http.StatusBadRequest, ErrCode: ErrCodeClientSideBadRequest, Message: "Bad Param"}
	ErrPanic         = DefaultError{statusCode: http.StatusInternalServerError, ErrCode: ErrCodeServerSideServerError, Message: "panic recovered"}
	ErrUnknownServer = DefaultError{statusCode: http.StatusInternalServerError, ErrCode: ErrCodeServerSideServerError, Message: "Unknown Server Error"}
)
