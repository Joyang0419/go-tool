package ginx_fx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"go-tool/pkg/ginx"
	"go-tool/pkg/ginx/ginx_api"
)

func TestRegisterRouter(t *testing.T) {
	c := ginx.Config{
		Port:            8000,
		ShutdownTimeout: 5 * time.Second,
		Mode:            gin.DebugMode,
	}
	_ = fx.New(
		RegisterRouter(NewTestController),
		Module(c),
	).Start(context.Background())

	// 發送測試請求
	resp, err := http.Get("http://localhost:8000/ping/456")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
}

type TestRouter struct{}

func NewTestController() *TestRouter {
	return &TestRouter{}
}

func (receiver *TestRouter) Routes(engine *gin.Engine) {
	ginx_api.New[pingSO, pingVO]().Method(http.MethodGet).Path("/ping/:userID").Service(Service).End(engine)
}

type pingSO struct {
	UserID string `uri:"userID" log:"true" binding:"required"`
}

type pingVO struct {
	UserID string `json:"userID"`
}

func Service(ctx context.Context, so pingSO) (vo pingVO, err error) {
	vo.UserID = so.UserID
	return vo, nil
}
