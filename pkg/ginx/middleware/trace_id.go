package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-tool/pkg/ginx/consts"
)

type TraceIDMiddleware struct {
}

func NewTraceIDMiddleware() *TraceIDMiddleware {
	return &TraceIDMiddleware{}
}

func (receiver *TraceIDMiddleware) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(consts.TraceIDKey)
		// 如果沒有，生成新的 traceId (使用 UUID)
		if traceID == "" {
			traceID = uuid.New().String()

			// 將生成的 traceId 設置到請求頭中
			c.Request.Header.Set(consts.TraceIDKey, traceID)
		}

		// 設置 traceId 到響應頭中
		c.Header(consts.TraceIDKey, traceID)

		// 設置到 gin context 中，方便後續使用
		c.Set(consts.TraceIDKey, traceID)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), consts.TraceIDKey, traceID))

		c.Next()
	}
}
