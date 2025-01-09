package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"framework/pkg/ginx/ginx_error"
	"framework/pkg/log"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 獲取堆疊跟踪
				stack := debug.Stack()

				// 解析堆疊資訊找到程式碼位置
				lines := strings.Split(string(stack), "\n")
				var errorLocation string

				for _, line := range lines {
					if strings.Contains(line, ".go:") &&
						!strings.Contains(line, "runtime/") &&
						!strings.Contains(line, "/gin") {
						if idx := strings.Index(line, " +0x"); idx != -1 {
							errorLocation = strings.TrimSpace(line[:idx])
						} else {
							errorLocation = strings.TrimSpace(line)
						}
						break
					}
				}

				// 轉換 error interface
				var webErr ginx_error.Error
				switch v := err.(type) {
				case ginx_error.Error:
					webErr = v
				case *ginx_error.Error:
					webErr = *v
				default:
					webErr = ginx_error.NewError(c.Request.Context(), http.StatusInternalServerError, ginx_error.ServerSideInternalErrCustomCode, fmt.Errorf("%+v", err))
				}

				scheme := c.GetHeader("X-Forwarded-Proto")
				if scheme == "" {
					if c.Request.TLS != nil {
						scheme = "https"
					} else {
						scheme = "http"
					}
				}
				// 印log
				log.Error(c.Request.Context(), "[RecoveryMiddleware]Panic recovered",
					zap.Any("error_location", errorLocation),
					zap.Any("request_method", c.Request.Method),
					zap.Any("request_url", scheme+"://"+c.Request.Host+c.Request.RequestURI),
					zap.Any("client_ip", c.ClientIP()),
					zap.Any("status_code", webErr.StatusCode),
					zap.Any("error_code", webErr.CustomCode),
					zap.Any("error_message", webErr.Message),
				)

				if webErr.StatusCode >= http.StatusBadRequest && webErr.StatusCode < http.StatusInternalServerError {
					// 回傳 4xx 錯誤
					c.AbortWithStatusJSON(webErr.StatusCode, webErr)
					return
				}

				// 回傳 500 錯誤
				c.AbortWithStatusJSON(http.StatusInternalServerError, webErr)
			}
		}()

		c.Next()
	}
}
