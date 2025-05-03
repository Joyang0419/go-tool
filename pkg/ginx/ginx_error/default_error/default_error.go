package default_error

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/cast"

	"go-tool/pkg/ginx/consts"
	"go-tool/pkg/ginx/ginx_error"
)

type DefaultError struct {
	StatusCode int         `json:"-"`
	CustomCode int         `json:"customCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	TraceID    string      `json:"traceID"`
}

func New(ctx context.Context, statusCode, customCode int, message ...string) DefaultError {
	return DefaultError{
		StatusCode: statusCode,
		CustomCode: customCode,
		Message:    lo.Ternary(len(message) > 0, message[0], ""),
		TraceID:    cast.ToString(ctx.Value(consts.TraceIDKey)),
	}
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
