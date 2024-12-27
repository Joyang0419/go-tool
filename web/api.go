package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tool/web/response"
)

type Application[REQ any, RESP any] func(ctx context.Context, request REQ) RESP

func ToHandlerFn[REQ any, RESP any](application Application[REQ, RESP]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request REQ
		if err := BindRequest(c, &request); err != nil {
			panic(fmt.Errorf("[ToHandlerFn]failed to bind request: %w", err))
		}

		result := application(c.Request.Context(), request)

		r := response.TResponse{
			Success: true,
			Data:    result,
		}

		c.JSON(http.StatusOK, r)
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
	if err := c.ShouldBindUri(request); err != nil {
		return fmt.Errorf("[BindRequest]failed to bind path parameters: %w", err)
	}
	if err := c.ShouldBindHeader(request); err != nil {
		return fmt.Errorf("[BindRequest]failed to bind header parameters: %w", err)
	}

	switch c.Request.Method {
	case http.MethodGet:
		// GET 請求只綁定 query parameters; 不綁定body
		return c.ShouldBindQuery(request)
	case http.MethodPatch, http.MethodPost, http.MethodPut:
		// POST, PUT, PATCH 請求綁定 body; 不綁定query parameters
		return c.ShouldBind(request)
	case http.MethodDelete:
		// DELETE 請求只綁定Path parameters; 不綁定body, query parameters
		return nil
	default:
		panic(fmt.Sprintf("[BindRequest]unsupported method: %s", c.Request.Method))
	}
}
