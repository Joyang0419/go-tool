package web

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	response2 "go-tool/web/response"
)

type TIRequest interface {
	Parse(c *gin.Context) error
}

type TApplication[REQ TIRequest, RESP any] func(ctx context.Context, request REQ) RESP

func ToHandlerFunc[REQ TIRequest, RESP any](application TApplication[REQ, RESP]) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqType := reflect.TypeOf((*REQ)(nil)).Elem()
		reqValue := reflect.New(reqType).Interface().(REQ)

		if err := reqValue.Parse(c); err != nil {
			panic(err)
		}

		response := response2.TResponse{
			Success: true,
			Data:    application(c.Request.Context(), reqValue),
		}

		c.JSON(http.StatusOK, response)
	}
}
