package ginx

import (
	"github.com/gin-gonic/gin"
)

type IRouter interface {
	// Routes 註冊路由, 例如: router.GET("/ping", func(c *gin.Context) {})
	Routes(engine *gin.Engine)
}

type IMiddleware interface {
	HandlerFunc() gin.HandlerFunc
}
