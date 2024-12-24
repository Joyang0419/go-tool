package web

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type TIController interface {
	// BasePath 路由前綴, 例如: /api/v1
	BasePath() string
	// Routes 註冊路由, 例如: router.GET("/ping", func(c *gin.Context) {})
	Routes(group *gin.RouterGroup)
	// Middlewares 註冊中間件, 例如: router.Use(middleware1, middleware2)
	Middlewares() []gin.HandlerFunc
}

func RegisterController(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(TIController)), fx.ResultTags(`group:"controllers"`)))
}

// TODO delete
type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (receiver *HealthController) BasePath() string {
	return "/health"
}

type HealthRequest struct{}

func (receiver HealthRequest) Parse(c *gin.Context) error {
	return nil
}

type HealthResponse struct{}

var HealthApplication = func(ctx context.Context, request HealthRequest) HealthResponse {
	return HealthResponse{}
}

func (receiver *HealthController) Routes(group *gin.RouterGroup) {
	group.GET("/")
}

func (receiver *HealthController) Middlewares() []gin.HandlerFunc {
	return nil
}
