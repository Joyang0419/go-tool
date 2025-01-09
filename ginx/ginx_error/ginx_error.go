package ginx_error

import (
	"context"

	"github.com/spf13/cast"

	"framework/pkg/ginx/consts"
)

type Error struct {
	StatusCode int         `json:"-"`
	CustomCode int         `json:"customCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	TraceID    string      `json:"traceId"`
}

const (
	ServerSideInternalErrCustomCode = 5000000

	ClientSideBadRequestCustomCode = 4000000
)

func (receiver *Error) Error() string {
	return receiver.Message
}

func NewError(ctx context.Context, statusCode, customCode int, err error, data ...any) Error {
	if len(data) > 0 {
		return Error{
			StatusCode: statusCode,
			CustomCode: customCode,
			Message:    err.Error(),
			TraceID:    cast.ToString(ctx.Value(consts.TraceIDKey)),
			Data:       data[0],
		}
	}
	return Error{
		StatusCode: statusCode,
		CustomCode: customCode,
		Message:    err.Error(),
		TraceID:    cast.ToString(ctx.Value(consts.TraceIDKey)),
	}
}
