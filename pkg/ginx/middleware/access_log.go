package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-tool/pkg/ginx/consts"
)

type AccessLogMiddleware struct{}

func NewAccessLogMiddleware() *AccessLogMiddleware {
	return &AccessLogMiddleware{}
}

func (receiver *AccessLogMiddleware) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 記錄請求資訊
		fields := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String(consts.TraceIDKey, c.GetString(consts.TraceIDKey)),
		}

		// GET 方法記錄 query parameters
		if c.Request.Method == http.MethodGet {
			fields = append(fields, slog.String("query", c.Request.URL.RawQuery))
		}

		// POST/PUT/PATCH 方法記錄 body
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch {
			var body []byte
			if c.Request.Body != nil {
				body, _ = io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
				if len(body) > 0 {
					bodyStr := string(body)
					bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
					fields = append(fields, slog.String("body", bodyStr))
				}
			}
		}

		slog.InfoContext(c.Request.Context(), "AccessLogMiddleware Access log", fields...)
		c.Next()
	}
}
