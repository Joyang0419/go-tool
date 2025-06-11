package resty

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"
	pkgerrors "github.com/pkg/errors"
	"github.com/samber/lo"

	"go-tool/pkg/httpclient"
	"go-tool/pkg/httpclient/consts"
	"go-tool/pkg/httpclient/logger"
)

type Client struct {
	client *resty.Client
}

func NewHTTPClient(config httpclient.Config) httpclient.IHTTPClient[resty.Client] {
	client := resty.New()
	client.SetBaseURL(config.BaseURL)
	if config.Logger != nil {
		logger.Init(config.Logger)
	}

	if config.EnableTrace {
		client.EnableTrace()
		client.OnBeforeRequest(EachRequestLog)
		client.OnAfterResponse(EachResponseLog)
	}

	client.OnError(OnErrorLog(config.BaseURL))

	return &Client{client: client}
}

func (receiver *Client) Request(ctx context.Context, param httpclient.RequestParam) (*http.Response, error) {
	if param.DisableLogResponseData {
		ctx = context.WithValue(ctx, consts.DisableLogRespData, struct{}{})
	}

	// 執行請求的核心邏輯
	executeRequest := func() (*http.Response, error) {
		req := receiver.client.NewRequest()
		req.SetContext(ctx)
		if param.Headers != nil {
			req.SetHeaders(param.Headers)
		}
		if param.PathParams != nil {
			req.SetPathParams(param.PathParams)
		}
		if param.QueryParams != nil {
			req.SetQueryParams(param.QueryParams)
		}
		if param.SuccessResponse != nil {
			req.SetResult(param.SuccessResponse)
		}
		if param.ErrorResponse != nil {
			req.SetError(param.ErrorResponse)
		}
		if param.Body != nil {
			req.SetBody(param.Body)
		}

		response, err := req.Execute(param.Method, param.URL)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"resty.Client.Request.error: Execute",
				slog.Any("error", err),
				slog.Any("param", param),
			)

			return nil, pkgerrors.Wrap(err, "resty.Client.Request.error: Execute")
		}

		if response.StatusCode() >= 300 {
			logger.Log.ErrorContext(
				ctx,
				"resty.Client.Request.error: status code not 2xx",
				slog.Int("statusCode", response.StatusCode()),
				slog.Any("param", param),
				slog.String("response", response.String()),
			)
			return nil, pkgerrors.Errorf(
				"resty.Client.Request.error: status code not 2xx: %d",
				response.StatusCode(),
			)
		}

		return response.RawResponse, nil
	}

	if lo.IsNil(param.RetryConfig) {
		return executeRequest()
	}

	return retry.DoWithData(
		executeRequest,
		param.RetryConfig.ToRetryOpts()...,
	)
}

func (receiver *Client) UploadFile(ctx context.Context, param httpclient.FileUploadParam) (*http.Response, error) {
	if param.DisableLogResponseData {
		ctx = context.WithValue(ctx, consts.DisableLogRespData, struct{}{})
	}

	executeRequest := func() (*http.Response, error) {
		req := receiver.client.NewRequest()
		req.SetContext(ctx)
		if param.Headers != nil {
			req.SetHeaders(param.Headers)
		}
		if param.PathParams != nil {
			req.SetPathParams(param.PathParams)
		}
		if param.QueryParams != nil {
			req.SetQueryParams(param.QueryParams)
		}
		if param.SuccessResponse != nil {
			req.SetResult(param.SuccessResponse)
		}
		if param.ErrorResponse != nil {
			req.SetError(param.ErrorResponse)
		}
		if param.FormData != nil {
			req.SetFormData(param.FormData)
		}

		// 設置文件
		for fieldName, fileInfo := range param.Files {
			if fileInfo.FilePath != "" && fileInfo.Reader != nil {
				return nil, pkgerrors.Errorf(
					"resty.Client.Request.error: field set both FilePath and Reader: %s",
					fieldName,
				)
			}

			if fileInfo.FilePath != "" {
				req.SetFile(fieldName, fileInfo.FilePath)
			}
			if fileInfo.Reader != nil {
				req.SetFileReader(fieldName, fileInfo.FileName, fileInfo.Reader)
			}
		}

		response, err := req.Execute(param.Method, param.URL)
		if err != nil {
			logger.Log.ErrorContext(
				ctx,
				"resty.Client.UploadFile.error: Execute",
				slog.Any("error", err),
				slog.Any("param", param),
			)
			return nil, pkgerrors.Wrap(err, "resty.Client.UploadFile.error: Execute")
		}

		return response.RawResponse, nil
	}

	if lo.IsNil(param.RetryConfig) {
		return executeRequest()
	}

	return retry.DoWithData(
		executeRequest,
		param.RetryConfig.ToRetryOpts()...,
	)
}

func (receiver *Client) Client() *resty.Client {
	return receiver.client
}
