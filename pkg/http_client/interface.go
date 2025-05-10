package http_client

import (
	"context"
	"net/http"
)

type IClient interface {
	Request(ctx context.Context, param RequestParam) (*http.Response, error)
}

type RequestParam struct {
	URL             string
	Method          string
	Headers         map[string]string
	PathParams      map[string]string
	QueryParams     map[string]string
	Body            any
	SuccessResponse any
	ErrorResponse   any
}
