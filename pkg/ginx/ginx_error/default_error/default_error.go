package default_error

import (
	"fmt"

	"go-tool/pkg/ginx/ginx_error"
)

type DefaultError struct {
	StatusCode int         `json:"-"`
	CustomCode int         `json:"customCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	TraceID    string      `json:"traceID"`
}

func (receiver DefaultError) Error() string {
	return fmt.Sprintf("ginx_error DefaultError status_code=%d custom_code=%d message=%s traceID=%s", receiver.StatusCode, receiver.CustomCode, receiver.Message, receiver.TraceID)
}

func (receiver DefaultError) HTTPStatusCode() int {
	return receiver.StatusCode
}

func (receiver DefaultError) Response() interface{} {

	return receiver
}

func (receiver DefaultError) SetTraceID(traceID string) ginx_error.IGinxError {
	receiver.TraceID = traceID

	return receiver
}

const (
	CustomCodeServerSideServerError = 500000
)

const (
	CustomCodeClientSideBadRequest = 400000
)
