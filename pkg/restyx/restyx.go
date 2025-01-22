package restyx

import (
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
)

// ClientConfig 定義 HTTP 客戶端的配置
type ClientConfig struct {
	// 超時相關配置
	Timeout struct {
		// Request 定義完整 HTTP 請求的超時時間，包含：
		// - DNS 解析
		// - 建立 TCP 連接
		// - TLS 握手
		// - 發送請求
		// - 服務器處理
		// - 接收響應
		// 建議值：
		// - 一般 API：10-30 秒
		// - 長時間操作：30-60 秒
		Request time.Duration `mapstructure:"request"`
	} `mapstructure:"timeout"`

	// 重試相關配置
	Retry struct {
		// Count 定義最大重試次數
		// 建議值：
		// - 一般操作：2-3 次
		// - 重要操作：3-5 次
		// 設置為 0 時不進行重試
		Count int `mapstructure:"count"`

		// WaitTime 定義重試之間的等待時間
		// 實際等待時間會根據重試策略調整
		// 建議值：100毫秒 - 1秒
		WaitTime time.Duration `mapstructure:"wait_time"`

		// MaxWaitTime 定義重試間最大等待時間
		// 用於避免指數退避時間過長
		// 建議值：WaitTime 的 2-3 倍
		MaxWaitTime time.Duration `mapstructure:"max_wait_time"`
	} `mapstructure:"retry"`
}

// NewClient creates a new HTTP client with the provided configuration
func NewClient(config ClientConfig) *resty.Client {
	client := resty.New()
	// Configure retries
	client.SetRetryCount(config.Retry.Count)
	client.SetRetryWaitTime(config.Retry.WaitTime)
	client.SetRetryMaxWaitTime(config.Retry.MaxWaitTime)

	return client
}

// InjectClient provides the HTTP client as a dependency using fx
func InjectClient(config ClientConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Provide(NewClient),
	)
}
