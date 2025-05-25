package resty

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-resty/resty/v2"
	pkgerrors "github.com/pkg/errors"

	"go-tool/pkg/httpclient"
	"go-tool/pkg/httpclient/consts"
)

type Client struct {
	client *resty.Client
}

func NewHTTPClient(config httpclient.Config) httpclient.IHTTPClient[*resty.Client] {
	client := resty.New()
	client.SetBaseURL(config.BaseURL)

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
		slog.ErrorContext(
			ctx,
			"resty.Client.Request.error: Execute",
			slog.Any("error", err),
			slog.Any("param", param),
		)

		return nil, pkgerrors.Wrap(err, "resty.Client.Request.error: Execute")
	}

	return response.RawResponse, nil
}

func (receiver *Client) UploadFile(ctx context.Context, param httpclient.FileUploadParam) (*http.Response, error) {
	if param.DisableLogResponseData {
		ctx = context.WithValue(ctx, consts.DisableLogRespData, struct{}{})
	}
	req := receiver.client.
		NewRequest().
		SetContext(ctx).
		SetHeaders(param.Headers).
		SetPathParams(param.PathParams).
		SetQueryParams(param.QueryParams).
		SetResult(param.SuccessResponse).
		SetError(param.ErrorResponse).
		SetFormData(param.FormData)

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
		slog.ErrorContext(
			ctx,
			"resty.Client.UploadFile.error: Execute",
			slog.Any("error", err),
			slog.Any("param", param),
		)
		return nil, pkgerrors.Wrap(err, "resty.Client.UploadFile.error: Execute")
	}

	return response.RawResponse, nil
}

func (receiver *Client) Client() *resty.Client {
	return receiver.client
}
