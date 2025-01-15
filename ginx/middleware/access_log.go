package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-tool/ginx/binding/logger"
	"go-tool/ginx/consts"
)

func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 記錄請求資訊
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("traceId", c.GetString(consts.TraceIDKey)),
		}

		// GET 方法記錄 query parameters
		if c.Request.Method == http.MethodGet {
			fields = append(fields, zap.String("query", c.Request.URL.RawQuery))
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
					fields = append(fields, zap.String("body", bodyStr))
				}
			}
		}

		logger.Info(c.Request.Context(), "[AccessLogMiddleware]Access log", fields...)
		c.Next()
	}
}
