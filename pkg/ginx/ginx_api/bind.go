package ginx_api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

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
		panic(fmt.Sprintf("ginx.BindRequest unsupported method: %s", c.Request.Method))
	}
}
