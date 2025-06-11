package ginx

import (
	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Routes(engine *gin.Engine)
}

type IMiddleware interface {
	HandlerFunc() gin.HandlerFunc
}
