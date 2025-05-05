package ginx_error

import (
	"context"
	"fmt"
	"net/http"

	"github.com/samber/lo"
	"github.com/spf13/cast"

	"go-tool/pkg/ginx/consts"
)

type DefaultError struct {
	StatusCode int         `json:"-"`
	CustomCode int         `json:"customCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	traceID    string      `json:"traceID"`
}

func New(ctx context.Context, statusCode, customCode int, message ...string) DefaultError {
	return DefaultError{
		StatusCode: statusCode,
		CustomCode: customCode,
		Message:    lo.Ternary(len(message) > 0, message[0], ""),
		traceID:    cast.ToString(ctx.Value(consts.TraceIDKey)),
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

func (receiver DefaultError) SetTraceID(traceID string) IGinxError {
	receiver.traceID = traceID

	return receiver
}

func (receiver DefaultError) TraceID() string {
	return receiver.traceID
}

const (
	CustomCodeServerSideServerError = 500000
)

const (
	CustomCodeClientSideBadRequest = 400000
)

var (
	ErrPanic    IGinxError = DefaultError{StatusCode: http.StatusInternalServerError, CustomCode: CustomCodeServerSideServerError}
	ErrParsedSO IGinxError = DefaultError{StatusCode: http.StatusBadRequest, CustomCode: CustomCodeClientSideBadRequest}
)
