package error

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"go-tool/pkg/ginx/consts"
)

var _ IGinxError = (*DefaultError)(nil)

type DefaultError struct {
	ErrCode int         `json:"errCode"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"traceID"`

	statusCode        int
	actualErr         error
	disableStackTrace bool
}

func (receiver DefaultError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("errorCode: %d\n", receiver.ErrCode))
	sb.WriteString(fmt.Sprintf("message: %s\n", receiver.Message))
	sb.WriteString(fmt.Sprintf("traceID: %s\n", receiver.TraceID))
	sb.WriteString(fmt.Sprintf("statusCode: %d\n", receiver.StatusCode()))
	if !receiver.disableStackTrace {
		sb.WriteString(fmt.Sprintf("stackTrace: %+v\n", receiver.actualErr))
	}

	return sb.String()
}

func (receiver DefaultError) StatusCode() int {
	if receiver.statusCode == 0 {
		return http.StatusInternalServerError
	}

	return receiver.statusCode
}

func (receiver DefaultError) SetStatusCode(statusCode int) IGinxError {
	receiver.statusCode = statusCode

	return receiver
}

func (receiver DefaultError) Response() interface{} {

	return receiver
}

func (receiver DefaultError) SetTraceID(traceID string) IGinxError {
	receiver.TraceID = traceID

	return receiver
}

func (receiver DefaultError) SetMessage(msg string) IGinxError {
	receiver.Message = msg

	return receiver
}

func (receiver DefaultError) SetErrorCode(errCode int) IGinxError {
	receiver.ErrCode = errCode

	return receiver
}

func (receiver DefaultError) SetActualError(err error) IGinxError {
	if receiver.actualErr != nil {
		receiver.actualErr = errors.Wrap(receiver.actualErr, err.Error())

		return receiver
	}

	receiver.actualErr = err
	return receiver
}

func (receiver DefaultError) SetDisableStackTrace() IGinxError {
	receiver.disableStackTrace = true

	return receiver
}

func (receiver DefaultError) SetCtx(ctx context.Context) IGinxError {
	receiver.TraceID = cast.ToString(ctx.Value(consts.TraceIDKey))

	return receiver
}
