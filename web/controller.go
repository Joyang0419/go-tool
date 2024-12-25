package web

import (
	"github.com/gin-gonic/gin"
)

type TIController interface {
	// BasePath 路由前綴, 例如: /api/v1
	BasePath() string
	// Routes 註冊路由, 例如: router.GET("/ping", func(c *gin.Context) {})
	Routes(group *gin.RouterGroup)
	// Middlewares 註冊中間件, 例如: router.Use(middleware1, middleware2)
	Middlewares() []gin.HandlerFunc
}
