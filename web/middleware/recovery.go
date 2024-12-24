package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tool/web/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 轉換 error interface
				var webErr response.TError
				switch v := err.(type) {
				case response.TError:
					webErr = v
				case *response.TError:
					webErr = *v
				default:
					webErr = response.NewError(http.StatusInternalServerError, 0, "Internal Server Error")
				}

				if webErr.StatusCode >= http.StatusBadRequest && webErr.StatusCode < http.StatusInternalServerError {
					// 回傳 4xx 錯誤
					c.AbortWithStatusJSON(
						webErr.StatusCode,
						response.TResponse{
							ErrorCode: webErr.CustomCode,
							Message:   webErr.Message,
						},
					)
					return

				}

				// 回傳 500 錯誤
				c.AbortWithStatusJSON(
					http.StatusInternalServerError, response.TResponse{
						ErrorCode: webErr.CustomCode,
						Message:   webErr.Message,
					},
				)
			}
		}()

		c.Next()
	}
}
