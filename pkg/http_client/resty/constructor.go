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

var _ http_client.IClient[*resty.Client] = (*Client)(nil)

func (receiver *Client) Request(ctx context.Context, param http_client.RequestParam) (*http.Response, error) {
	request := receiver.client.NewRequest().
		SetContext(ctx).
		SetHeaders(param.Headers).
		SetPathParams(param.PathParams).
		SetQueryParams(param.QueryParams).
		SetResult(param.SuccessResponse).
		SetError(param.ErrorResponse)

	restyResp, err := request.Execute(param.Method, param.URL)
	if err != nil {
		slog.ErrorContext(ctx, "resty.Client.NewRequest error",
			slog.Any("error", err.Error()),
			slog.Any("param", param),
		)

		return nil, pkgerrors.Wrap(err, "resty.Client.NewRequest error")
	}

	return restyResp.RawResponse, nil
}

func (receiver *Client) GetClient() *resty.Client {
	return receiver.client
}
