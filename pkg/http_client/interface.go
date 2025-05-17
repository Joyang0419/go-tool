package http_client

import (
	"context"
	"net/http"

	"github.com/avast/retry-go/v4"
)

type IClient[TClient any] interface {
	Request(ctx context.Context, param RequestParam, retry ...rerty.) (*http.Response, error)
	// GetClient 其實沒有想開這個;
	// 但避免Request, 寫得不夠「通用」，還是開一個可以直接使用Lib的拿到原生的Client去給使用者操作
	GetClient() TClient
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

func qq() {
	retry.Do(func() error {

	})

	retry.IsRecoverable()
}
