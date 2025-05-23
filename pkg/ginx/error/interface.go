package error

type IGinxError interface {
	Error() string
	SetStatusCode(int) IGinxError
	StatusCode() int
	Response() interface{}
	SetTraceID(traceID string) IGinxError
	SetActualError(err error) IGinxError
	SetDisableStackTrace() IGinxError
	SetMessage(msg string) IGinxError
}
