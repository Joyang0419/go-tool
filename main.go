package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"go-tool/convert"
	"go-tool/web"
)

func main() {
	fx.New(
		web.RegisterController(NewHealthController),
		web.InjectServer(web.NewServerConfig(8081, 5*time.Second)),
	).Run()
}

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (receiver *HealthController) BasePath() string {
	return "/health"
}

type HealthRequest struct {
	ID int `json:"id"`
}

func (receiver HealthRequest) Parse(c *gin.Context) (web.IRequest[any], error) {
	receiver.ID = convert.ToInt(c.Param("id"))

	return receiver, nil
}

type HealthResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

var HealthApplication web.Application[HealthRequest, HealthResponse] = func(ctx context.Context, request HealthRequest) HealthResponse {
	return HealthResponse{
		Message: "pong",
		ID:      request.ID,
	}
}

func (receiver *HealthController) Routes(group *gin.RouterGroup) {
	group.GET("/:id", web.ToHandlerFn(HealthApplication))
}

func (receiver *HealthController) Middlewares() []gin.HandlerFunc {
	return nil
}
