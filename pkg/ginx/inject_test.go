package ginx

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"go-tool/pkg/ginx/ginx_error"
)

func TestRegisterRouter(t *testing.T) {
	_ = fx.New(
		RegisterRouter(NewTestController),
		InjectServer(ServerConfig{
			Port:            8000,
			ShutdownTimeout: 5 * time.Second,
		}),
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
	engine.GET("/ping/:userId", ToHandlerFn(receiver.ping))
}

type pingRequest struct {
	UserId string `uri:"userId"`
}

type pingResponse struct {
	UserId string `json:"userId"`
}

// 使用者的input, 就是所有會用到 例如來自 postBody, queryParams, or uri 的資料; 全部定義在request struct裡面
func (receiver *TestRouter) ping(ctx context.Context, request pingRequest) (response pingResponse) {
	response.UserId = request.UserId
	slog.Info("ping", slog.String("userId", request.UserId))
	// 順手測panic機制: 會被recovery 接到, 並且回傳500, 然後我要看會印traceId嗎? 有的話就是我要的
	panic(ginx_error.NewError(ctx, http.StatusInternalServerError, ginx_error.ServerSideInternalErrCustomCode, fmt.Errorf("test panic")))
	return response
}
