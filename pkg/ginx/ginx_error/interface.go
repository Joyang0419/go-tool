package ginx_error

type IGinxError interface {
	Error() string
	HTTPStatusCode() int
	Response() interface{}
	SetTraceID(traceID string) IGinxError
	TraceID() string
}
