package resty

import (
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"

	"go-tool/pkg/httpclient/consts"
	"go-tool/pkg/httpclient/logger"
)

func EachRequestLog(client *resty.Client, request *resty.Request) error {
	_ = client
	// 取得完整 URL
	fullURL := fmt.Sprintf("%s%s", client.BaseURL, request.URL)
	logger.Log.InfoContext(
		request.Context(),
		"resty.EachRequestLog",
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
		slog.Any("requestBody", response.Request.Body),
		slog.Any("queryParam", response.Request.QueryParam),
		slog.Any("pathParam", response.Request.PathParams),
		slog.Any("formData", response.Request.FormData),
		slog.Float64("responseTime", response.Time().Seconds()),
	}
	if traceInfo.TotalTime > 0 {
		slogAttrs = append(
			slogAttrs,
			slog.Float64("dnsLookup", traceInfo.DNSLookup.Seconds()),
			slog.Float64("connTime", traceInfo.ConnTime.Seconds()),
			slog.Float64("tcpConnTime", traceInfo.TCPConnTime.Seconds()),
			slog.Float64("tlsHandshake", traceInfo.TLSHandshake.Seconds()),
			slog.Float64("serverTime", traceInfo.ServerTime.Seconds()),
			slog.Float64("responseTime", traceInfo.ResponseTime.Seconds()),
			slog.Float64("totalTime", traceInfo.TotalTime.Seconds()),
			slog.Float64("connIdleTime", traceInfo.ConnIdleTime.Seconds()),
			slog.Bool("isConnReused", traceInfo.IsConnReused),
			slog.Bool("isConnWasIdle", traceInfo.IsConnWasIdle),
			slog.Int("requestAttempt", traceInfo.RequestAttempt),
			slog.String("remoteAddr", traceInfo.RemoteAddr.String()),
		)
	}

	isNeedLogResponseData := lo.IsNil(ctx.Value(consts.DisableLogRespData))
	if isNeedLogResponseData {
		slogAttrs = append(
			slogAttrs,
			slog.String("responseBody", string(response.Body())),
		)
	}

	if response.StatusCode() >= 300 {
		logger.Log.ErrorContext(
			ctx,
			"resty.EachResponseLog.error",
			slogAttrs...,
		)

		return nil
	}

	logger.Log.InfoContext(
		ctx,
		"resty.EachResponseLog",
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
		logger.Log.ErrorContext(
			request.Context(),
			"resty.OnErrorLog.error",
			slog.Any("error", err),
			slog.String("url", fullURL),
			slog.String("method", request.Method),
			slog.Any("header", request.Header),
			slog.Any("requestBody", request.Body),
			slog.Any("queryParam", request.QueryParam),
			slog.Any("pathParam", request.PathParams),
			slog.Any("formData", request.FormData),
		)
	}
}
