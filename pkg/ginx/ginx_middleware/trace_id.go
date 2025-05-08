package ginx_middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-tool/pkg/ginx/ginx_consts"
)

// TraceID 從 gin context 中獲取 traceId, 如果沒有，生成新的 traceId (使用 UUID)
func TraceID(c *gin.Context) {
	traceID := c.GetHeader(ginx_consts.TraceIDKey)
	// 2. 如果沒有，生成新的 traceId (使用 UUID)
	if traceID == "" {
		traceID = uuid.New().String()

		// 將生成的 traceId 設置到請求頭中
		c.Request.Header.Set(ginx_consts.TraceIDKey, traceID)
	}

	// 3. 設置 traceId 到響應頭中
	c.Header(ginx_consts.TraceIDKey, traceID)

	// 4. 設置到 gin context 中，方便後續使用
	c.Set(ginx_consts.TraceIDKey, traceID)
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), ginx_consts.TraceIDKey, traceID))
}

func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		TraceID(c)
		c.Next()
	}
}
