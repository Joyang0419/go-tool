package ginx_api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tool/pkg/ginx/ginx_error"
)

type THandler func(ctx context.Context, so any) (vo any, err error)

type API[TSO, TVO any] struct {
	httpMethod     string
	path           string
	handler        THandler
	beforeHandlers []gin.HandlerFunc
	afterHandlers  []gin.HandlerFunc
}

func New[TSO, TVO any]() *API[TSO, TVO] {
	return &API[TSO, TVO]{}
}

func (receiver *API[TSO, TVO]) SetHTTPMethod(httpMethod string) *API[TSO, TVO] {
	receiver.httpMethod = httpMethod

	return receiver
}

func (receiver *API[TSO, TVO]) SetPath(path string) *API[TSO, TVO] {
	receiver.path = path

	return receiver
}

func (receiver *API[TSO, TVO]) SetHandler(h THandler) *API[TSO, TVO] {
	receiver.handler = h

	return receiver
}

func (receiver *API[TSO, TVO]) SetBeforeHandler(handlers ...gin.HandlerFunc) *API[TSO, TVO] {
	receiver.beforeHandlers = append(receiver.beforeHandlers, handlers...)

	return receiver
}

func (receiver *API[TSO, TVO]) SetAfterHandler(handlers ...gin.HandlerFunc) *API[TSO, TVO] {
	receiver.afterHandlers = append(receiver.afterHandlers, handlers...)

	return receiver
}

func (receiver *API[TSO, TVO]) Register(engine *gin.Engine) {
	mainHandlerFn := func(c *gin.Context) {
		var so TSO
		if err := BindRequest(c, &so); err != nil {
			slog.ErrorContext(
				c.Request.Context(), "API BindRequest error",
				slog.Any("error", err),
				slog.String("method", receiver.httpMethod),
				slog.String("path", receiver.path),
			)

			_ = c.Error(ginx_error.ErrorInvalidRequest(c.Request.Context(), err))
			return
		}

		vo, err := receiver.handler(c.Request.Context(), so)
		if err != nil {
			slog.ErrorContext(
				c.Request.Context(), "API Handler error",
				slog.Any("error", err),
				slog.String("method", receiver.httpMethod),
				slog.String("path", receiver.path),
			)

			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, vo)
	}

	handlerFns := append(receiver.beforeHandlers, mainHandlerFn)
	handlerFns = append(handlerFns, receiver.afterHandlers...)

	engine.Handle(receiver.httpMethod, receiver.path, handlerFns...)
}
