package ginx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"

	"go-tool/pkg/ginx/consts"
	"go-tool/pkg/ginx/ginx_error"
)

type API[REQ any, RESP any] func(ctx context.Context, request REQ) (RESP, error)

func ToHandlerFn[REQ any, RESP any](api API[REQ, RESP]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request REQ
		if err := BindRequest(c, &request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ginx_error.NewError(c.Request.Context(), http.StatusBadRequest, ginx_error.ClientSideBadRequestCustomCode, err))
			return
		}

		result, err := api(c.Request.Context(), request)
		if err != nil {
			var ginErr ginx_error.Error
			if errors.As(err, &ginErr) {
				if ginErr.StatusCode >= http.StatusInternalServerError {
					scheme := c.GetHeader("X-Forwarded-Proto")
					if scheme == "" {
						if c.Request.TLS != nil {
							scheme = "https"
						} else {
							scheme = "http"
						}
					}
					pc, _, _, _ := runtime.Caller(1)
					errorLocation := runtime.FuncForPC(pc).Name()
					slog.ErrorContext(c.Request.Context(), "[ToHandlerFn]server error",
						slog.Any("error_location", errorLocation),
						slog.Any("request_method", c.Request.Method),
						slog.Any("request_url", scheme+"://"+c.Request.Host+c.Request.RequestURI),
						slog.Any("client_ip", c.ClientIP()),
						slog.Any("status_code", ginErr.StatusCode),
						slog.Any("error_code", ginErr.CustomCode),
						slog.Any("error_message", ginErr.Message),
						slog.Any(consts.TraceIDKey, c.GetString(consts.TraceIDKey)),
					)
				}

				c.AbortWithStatusJSON(ginErr.StatusCode, ginErr)
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, ginx_error.NewError(c.Request.Context(), http.StatusInternalServerError, ginx_error.ServerSideInternalErrCustomCode, err))
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// BindRequest 綁定請求參數
/* example struct:
type UserRequest struct {
    ID      int    `uri:"id"`        // path parameter /:id
	Message string `query:"message"` // query parameter ?message=hello
    Tags        []string `form:"tags" json:"tags"`        // form-data array or json array
	file 	  *multipart.FileHeader `form:"file"` // form-data file
	lang 	string `header:"Accept-Language"` // header parameter Accept-Language
}
*/
func BindRequest(c *gin.Context, request interface{}) error {
	// 為何不判斷err, 因為每個should都會去判定validate, 如果用json的，前面還沒吃到值，就會validate err, 所以先不判斷err, 最後在統一判斷
	_ = c.ShouldBindUri(request)
	_ = c.ShouldBindHeader(request)

	switch c.Request.Method {
	case http.MethodGet:
		// GET 請求只綁定 query parameters; 不綁定body
		if len(c.Request.URL.Query()) == 0 {
			return binding.Validator.ValidateStruct(request)
		}
		return c.ShouldBindQuery(request)
	case http.MethodPatch, http.MethodPost, http.MethodPut:
		// POST, PUT, PATCH 請求綁定 body; 不綁定query parameters
		// 檢查請求體是否為空
		if c.Request.Body == nil || c.Request.ContentLength == 0 {
			return binding.Validator.ValidateStruct(request)
		}

		return c.ShouldBind(request)
	case http.MethodDelete:
		// DELETE 請求只綁定Path parameters; 不綁定body, query parameters
		return binding.Validator.ValidateStruct(request)
	default:
		panic(fmt.Sprintf("[BindRequest]unsupported method: %s", c.Request.Method))
	}
}
