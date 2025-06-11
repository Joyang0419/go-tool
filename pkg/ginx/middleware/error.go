package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"go-tool/pkg/ginx/error"
)

type ErrorMiddleware struct{}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (receiver *ErrorMiddleware) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// 記錄所有錯誤
		for _, ginErr := range c.Errors {
			slog.ErrorContext(c.Request.Context(), "ErrorMiddleware error",
				slog.Any("error", ginErr.Err),
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.String("client_ip", c.ClientIP()),
			)
		}

		err := c.Errors[0].Err
		var errGinx error.IGinxError
		if errors.As(err, &errGinx) {
			c.AbortWithStatusJSON(errGinx.StatusCode(), errGinx.Response())
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
}
