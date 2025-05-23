package resty

import (
	"log/slog"

	"github.com/go-resty/resty/v2"
)

func EachRequestLog(c *resty.Client, req *resty.Request) error {
	return nil
}

func EachResponseLog(c *resty.Client, resp *resty.Response) error {
	return nil
}

func OnErrorLog(c *resty.Request, err error) {
	slog.ErrorContext(
		c.Context(),
		"resty.OnErrorLog error",
		slog.Any("error", err),
		slog.String("url", c.URL),
		slog.String("method", c.Method),
		slog.Any("headers", c.Header),
		slog.Any("body", c.Body),
		slog.Any("queryParams", c.QueryParam),
		slog.Any("pathParams", c.PathParams),
		slog.Any("formData", c.FormData),
		slog.Any("result", c.Result),
	)
}
