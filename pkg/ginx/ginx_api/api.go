package ginx_api

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
	"go-tool/pkg/ginx/ginx_error"
)

type TService[TSO any, TVO any] func(ctx context.Context, so TSO) (vo TVO, err error)

type API[TSO, TVO any] struct {
	httpMethod     string
	path           string
	service        TService[TSO, TVO]
	beforeHandlers []gin.HandlerFunc
	afterHandlers  []gin.HandlerFunc
}

func New[TSO, TVO any]() *API[TSO, TVO] {
	return &API[TSO, TVO]{}
}

func (receiver *API[TSO, TVO]) HTTPMethod(httpMethod string) *API[TSO, TVO] {
	receiver.httpMethod = httpMethod

	return receiver
}

func (receiver *API[TSO, TVO]) Path(path string) *API[TSO, TVO] {
	receiver.path = path

	return receiver
}

func (receiver *API[TSO, TVO]) Service(service TService[TSO, TVO]) *API[TSO, TVO] {
	receiver.service = service

	return receiver
}

func (receiver *API[TSO, TVO]) BeforeHandlers(handlers ...gin.HandlerFunc) *API[TSO, TVO] {
	receiver.beforeHandlers = append(receiver.beforeHandlers, handlers...)

	return receiver
}

func (receiver *API[TSO, TVO]) AfterHandlers(handlers ...gin.HandlerFunc) *API[TSO, TVO] {
	receiver.afterHandlers = append(receiver.afterHandlers, handlers...)

	return receiver
}

func (receiver *API[TSO, TVO]) End(engine *gin.Engine) {
	if lo.IsNil(receiver.service) {
		panic(fmt.Sprintf("API service is nil, path: %s, method: %s", receiver.path, receiver.httpMethod))
	}

	mainHandlerFn := func(c *gin.Context) {
		var so TSO
		if err := BindRequest(c, &so); err != nil {
			slog.ErrorContext(
				c.Request.Context(), "API BindRequest error",
				slog.Any("error", err),
				slog.String("method", receiver.httpMethod),
				slog.String("path", receiver.path),
			)

			_ = c.Error(ginx_error.ErrParsedSO.SetTraceID(cast.ToString(c.Request.Context().Value(consts.TraceIDKey))))
			return
		}

		// 使用 log 標籤來決定要記錄哪些字段
		loggableFields := getLoggableFields(so)
		if len(loggableFields) > 0 {
			slogAttrs := lo.MapToSlice(loggableFields, func(key string, value any) any {
				return slog.Any(key, value)
			})

			slog.InfoContext(c.Request.Context(), "so fields log",
				slogAttrs...,
			)
		}

		vo, err := receiver.service(c.Request.Context(), so)
		if err != nil {
			slog.ErrorContext(
				c.Request.Context(), "API Service error",
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
