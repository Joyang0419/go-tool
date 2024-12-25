package web

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"go-tool/web/response"
)

type TIRequest[T any] interface {
	Parse(c *gin.Context) (TIRequest[T], error)
}

type TApplication[REQ any, RESP any] func(ctx context.Context, request REQ) RESP

func ToHandlerFn[REQ any, RESP any](application TApplication[REQ, RESP]) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查類型並建立實例
		reqType := reflect.TypeOf((*REQ)(nil)).Elem()
		if reqType == nil {
			panic(response.NewError(http.StatusInternalServerError, 0, "[ToHandlerFn]reqType == nil"))
		}

		// 建立新實例
		reqInstance := reflect.New(reqType)
		if !reqInstance.IsValid() {
			panic(response.NewError(http.StatusInternalServerError, 0, "[ToHandlerFn]reqInstance.IsValid() == false"))
		}

		// 嘗試轉換為 TIRequest 接口
		reqValue, ok := reqInstance.Interface().(TIRequest[any])
		if !ok {
			panic(response.NewError(http.StatusInternalServerError, 0, "[ToHandlerFn]reqValue.(TIRequest) == false"))
		}

		// 解析請求
		parsedREQ, err := reqValue.Parse(c)
		if err != nil {
			panic(response.NewError(http.StatusBadRequest, 0, fmt.Sprintf("[ToHandlerFn]reqValue.Parse(c) error: %v", err)))
		}

		// 檢查 parsedREQ 是否為 nil
		if parsedREQ == nil {
			panic(response.NewError(http.StatusInternalServerError, 0, "[ToHandlerFn]parsedREQ == nil"))
		}

		// 嘗試轉換回 REQ 類型
		finalReq, ok := parsedREQ.(REQ)
		if !ok {
			panic(response.NewError(http.StatusInternalServerError, 0, "[ToHandlerFn]parsedREQ.(REQ) == false"))
		}

		// 執行應用邏輯
		result := application(c.Request.Context(), finalReq)

		r := response.TResponse{
			Success: true,
			Data:    result,
		}

		c.JSON(http.StatusOK, r)
	}
}
