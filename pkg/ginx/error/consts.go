package error

import (
	"net/http"
)

const (
	ErrCodeServerSideServerError = 500000
)

const (
	ErrCodeClientSideBadRequest = 400000

	ErrCodeClientSideUnauthorized = 401000
	ErrCodeClientSideTokenEmpty   = 401001 // token 為空
	ErrCodeClientSideTokenInvalid = 401002 // token 驗證失敗
)

func ErrBadParam() IGinxError {
	return DefaultError{
		statusCode: http.StatusBadRequest,
		ErrCode:    ErrCodeClientSideBadRequest,
		Message:    "Bad Param",
	}
}

func ErrPanic() IGinxError {
	return DefaultError{
		statusCode: http.StatusInternalServerError,
		ErrCode:    ErrCodeServerSideServerError,
		Message:    "panic recovered",
	}
}

func ErrUnknownServer() IGinxError {
	return DefaultError{
		statusCode: http.StatusInternalServerError,
		ErrCode:    ErrCodeServerSideServerError,
		Message:    "Unknown Server Error",
	}
}

func ErrUnauthorized() IGinxError {
	return DefaultError{
		statusCode: http.StatusUnauthorized,
		ErrCode:    ErrCodeClientSideUnauthorized,
		Message:    "Unauthorized",
	}
}

func ErrTokenEmpty() IGinxError {
	return ErrUnauthorized().
		SetErrorCode(ErrCodeClientSideTokenEmpty).
		SetMessage("token is required")
}

func ErrTokenInvalid() IGinxError {
	return ErrUnauthorized().
		SetErrorCode(ErrCodeClientSideTokenInvalid).
		SetMessage("token invalid")
}
