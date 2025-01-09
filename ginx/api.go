package ginx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type API[REQ any, RESP any] func(ctx context.Context, request REQ) RESP

func ToHandlerFn[REQ any, RESP any](api API[REQ, RESP]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request REQ
		if err := BindRequest(c, &request); err != nil {
			panic(fmt.Errorf("[ToHandlerFn]failed to bind request: %w", err))
		}

		result := api(c.Request.Context(), request)

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
