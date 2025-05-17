package resty

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-resty/resty/v2"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"go-tool/pkg/http_client"
)

type Client struct {
	client *resty.Client
}

func NewRestyClient(config Config) *Client {
	client := resty.New()
	client.SetBaseURL(config.BaseURL)

	if config.EnableAfterResponseLog {
		client.OnBeforeRequest(EachRequestLog)
	}
	if config.EnableAfterResponseLog {
		client.OnAfterResponse(EachResponseLog)
	}

	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return slices.Contains(config.RetryStatusCodes, r.StatusCode())
		},
	)

	client.OnError(OnErrorLog)

	return &Client{client: client}
}

var _ http_client.IClient[*resty.Client] = (*Client)(nil)

func (receiver *Client) Request(ctx context.Context, param http_client.RequestParam) (*http.Response, error) {
	response, err := receiver.client.
		NewRequest().
		SetContext(ctx).
		SetHeaders(param.Headers).
		SetPathParams(param.PathParams).
		SetQueryParams(param.QueryParams).
		SetResult(param.SuccessResponse).
		SetError(param.ErrorResponse).
		Execute(param.Method, param.URL)
	if err != nil {
		slog.ErrorContext(ctx, "resty.Client.Request Execute error",
			slog.Any("error", err),
			slog.Any("param", param),
		)

		return nil, pkgerrors.Wrap(err, "resty.Client.Request Execute error")
	}

	return response.RawResponse, nil
}

func (receiver *Client) GetClient() *resty.Client {
	return receiver.client
}
