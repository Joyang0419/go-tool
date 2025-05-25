package resty

import (
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"

	"go-tool/pkg/httpclient/consts"
)

func EachRequestLog(client *resty.Client, request *resty.Request) error {
	_ = client
	// 取得完整 URL
	fullURL := fmt.Sprintf("%s%s", client.BaseURL, request.URL)
	slog.InfoContext(
		request.Context(),
		"resty.EachRequestLog.info",
		slog.String("url", fullURL),
		slog.String("method", request.Method),
		slog.Any("header", request.Header),
		slog.Any("pathParam", request.PathParams),
		slog.Any("queryParam", request.QueryParam),
		slog.Any("body", request.Body),
		slog.Any("formData", request.FormData),
	)

	return nil
}

func EachResponseLog(client *resty.Client, response *resty.Response) error {
	traceInfo := response.Request.TraceInfo()
	ctx := response.Request.Context()
	fullURL := fmt.Sprintf("%s%s", client.BaseURL, response.Request.URL)
	slogAttrs := []any{
		slog.Int("statusCode", response.StatusCode()),
		slog.String("url", fullURL),
		slog.String("method", response.Request.Method),
		slog.Any("header", response.Request.Header),
		slog.Any("body", response.Request.Body),
		slog.Any("queryParam", response.Request.QueryParam),
		slog.Any("pathParam", response.Request.PathParams),
		slog.Any("formData", response.Request.FormData),
		slog.Duration("responseTime", response.Time()),
	}
	if traceInfo.TotalTime > 0 {
		slogAttrs = append(slogAttrs, slog.Any("traceInfo", traceInfo))
	}

	isNeedLogResponseData := lo.IsNil(ctx.Value(consts.DisableLogRespData))
	if isNeedLogResponseData {
		slogAttrs = append(
			slogAttrs,
			slog.String("body", string(response.Body())),
		)
	}

	if response.StatusCode() >= 300 {
		slog.ErrorContext(
			ctx,
			"resty.EachResponseLog.error",
			slogAttrs...,
		)

		return nil
	}

	slog.InfoContext(
		ctx,
		"resty.EachResponseLog.info",
		slogAttrs...,
	)

	return nil
}

/*
OnErrorLog 在 Resty 请求过程中发生错误时被调用的回调函数

调用时机：
  - 当 HTTP 请求执行过程中遇到错误时（如连接失败、超时等网络错误）
  - 在请求被发送后但未能成功完成时
  - 在 Resty 客户端的错误处理链中作为错误处理器被触发

该函数被设计用作 client.SetErrorHook() 的回调函数，用于详细记录请求失败的相关信息。
*/
func OnErrorLog(baseURL string) resty.ErrorHook {
	return func(request *resty.Request, err error) {
		fullURL := fmt.Sprintf("%s%s", baseURL, request.URL)
		slog.ErrorContext(
			request.Context(),
			"resty.OnErrorLog.error",
			slog.Any("error", err),
			slog.String("url", fullURL),
			slog.String("method", request.Method),
			slog.Any("header", request.Header),
			slog.Any("body", request.Body),
			slog.Any("queryParam", request.QueryParam),
			slog.Any("pathParam", request.PathParams),
			slog.Any("formData", request.FormData),
		)
	}
}
