package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/spf13/cast"

	"go-tool/pkg/ginx/consts"
	ginxerror "go-tool/pkg/ginx/error"
)

type TService[TSO any, TVO any] func(ctx context.Context, so TSO) (vo TVO, err error)

type TAPI[TSO, TVO any] struct {
	method         string
	path           string
	service        TService[TSO, TVO]
	beforeHandlers []gin.HandlerFunc
	afterHandlers  []gin.HandlerFunc
}

func New[TSO, TVO any]() *TAPI[TSO, TVO] {
	return &TAPI[TSO, TVO]{}
}

func (receiver *TAPI[TSO, TVO]) Method(method string) *TAPI[TSO, TVO] {
	receiver.method = method

	return receiver
}

func (receiver *TAPI[TSO, TVO]) Path(path string) *TAPI[TSO, TVO] {
	receiver.path = path

	return receiver
}

func (receiver *TAPI[TSO, TVO]) Service(service TService[TSO, TVO]) *TAPI[TSO, TVO] {
	receiver.service = service

	return receiver
}

func (receiver *TAPI[TSO, TVO]) BeforeHandlers(handlers ...gin.HandlerFunc) *TAPI[TSO, TVO] {
	receiver.beforeHandlers = append(receiver.beforeHandlers, handlers...)

	return receiver
}

func (receiver *TAPI[TSO, TVO]) AfterHandlers(handlers ...gin.HandlerFunc) *TAPI[TSO, TVO] {
	receiver.afterHandlers = append(receiver.afterHandlers, handlers...)

	return receiver
}

func (receiver *TAPI[TSO, TVO]) End(engine *gin.Engine) {
	if lo.IsNil(receiver.service) {
		panicMsg := fmt.Sprintf("TAPI.End service is nil, path: %s, method: %s", receiver.path, receiver.method)
		panic(panicMsg)
	}

	mainHandlerFn := func(c *gin.Context) {
		var so TSO
		if err := BindRequest(c, &so); err != nil {
			slog.ErrorContext(
				c.Request.Context(), "TAPI.mainHandlerFn.error: BindRequest",
				slog.Any("error", err),
				slog.String("method", receiver.method),
				slog.String("path", receiver.path),
			)

			_ = c.Error(ginxerror.ErrBadParam.SetTraceID(cast.ToString(c.Request.Context().Value(consts.TraceIDKey))))
			return
		}

		// 使用 log 標籤來決定要記錄哪些字段
		loggableFields := getLoggableFields(so)
		if len(loggableFields) > 0 {
			slogAttrs := lo.MapToSlice(
				loggableFields,
				func(key string, value any) any {
					return slog.Any(key, value)
				},
			)

			slog.InfoContext(
				c.Request.Context(),
				"TAPI.mainHandlerFn.SO.loggableFields",
				slogAttrs...,
			)
		}

		vo, err := receiver.service(c.Request.Context(), so)
		if err != nil {
			slog.ErrorContext(
				c.Request.Context(), "TAPI.mainHandlerFn.error: receiver.service",
				slog.Any("error", err),
				slog.String("method", receiver.method),
				slog.String("path", receiver.path),
			)

			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, vo)
	}

	handlerFns := append(receiver.beforeHandlers, mainHandlerFn)
	handlerFns = append(handlerFns, receiver.afterHandlers...)

	engine.Handle(receiver.method, receiver.path, handlerFns...)
}

// 通用函數來提取帶有 log:"true" 標籤的字段
func getLoggableFields(v interface{}) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if logTag := field.Tag.Get("log"); logTag == "true" {
			fieldValue := val.Field(i)

			// 使用 JSON tag 作為 key（如果存在），否則使用字段名
			key := field.Name
			if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				// 處理 json tag 中的逗號（例如 `json:"name,omitempty"`）
				if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
					jsonTag = jsonTag[:commaIdx]
				}
				key = jsonTag
			}

			result[key] = fieldValue.Interface()
		}
	}

	return result
}
