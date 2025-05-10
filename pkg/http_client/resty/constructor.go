package resty

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-resty/resty/v2"
	pkgerrors "github.com/pkg/errors"

	"go-tool/pkg/http_client"
)

type Client struct {
	client *resty.Client
}

func NewRestyClient(config Config) *Client {
	client := resty.New()
	client.SetBaseURL(config.BaseURL)

	if config.EnableBeforeRequestLog {
		client.OnBeforeRequest(EachRequestLog)
	}
	if config.EnableAfterResponseLog {
		client.OnAfterResponse(EachResponseLog)
	}

	client.OnError(OnErrorLog)

	return &Client{client: client}
}

var _ http_client.IClient = (*Client)(nil)

func (receiver *Client) Request(ctx context.Context, param http_client.RequestParam) (*http.Response, error) {
	restyResp, err := receiver.client.NewRequest().
		SetContext(ctx).
		SetHeaders(param.Headers).
		SetPathParams(param.PathParams).
		SetQueryParams(param.QueryParams).
		SetBody(param.Body).
		SetResult(param.SuccessResponse).
		SetError(param.ErrorResponse).
		Execute(param.Method, param.URL)
	if err != nil {
		slog.ErrorContext(ctx, "resty.Client.NewRequest error",
			slog.Any("error", err.Error()),
			slog.Any("param", param),
		)

		return nil, pkgerrors.Wrap(err, "resty.Client.NewRequest error")
	}

	return restyResp.RawResponse, nil
}
