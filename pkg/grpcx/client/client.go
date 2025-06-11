package client

import (
	"context"
	"crypto/tls"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"`
	URL     string        `mapstructure:"url" yaml:"url"`
	UseTLS  bool          `mapstructure:"use_tls" yaml:"use_tls"`
}

func NewClient(config Config, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	var options []grpc.DialOption

	// TLS 設定
	if config.UseTLS {
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 添加其他選項
	options = append(options, opt...)
	options = append(options,
		grpc.WithConnectParams(
			grpc.ConnectParams{
				// MinConnectTimeout: 每次連接嘗試的最小超時時間
				// 如果連接在此時間內未建立，會觸發重試機制
				MinConnectTimeout: 5 * time.Second,

				Backoff: backoff.Config{
					// BaseDelay: 第一次重試前的等待時間
					BaseDelay: 1.0 * time.Second,

					// Multiplier: 每次重試失敗後，下次等待時間的乘數
					// 例如：第一次等待1秒，第二次1.6秒，第三次2.56秒...
					Multiplier: 1.6,

					// Jitter: 在重試間隔時間上添加隨機波動
					// 避免多個客戶端同時重試造成服務器壓力
					// 例如：如果間隔是1秒，Jitter為0.2，實際等待時間會在0.8-1.2秒之間
					Jitter: 0.2,

					// MaxDelay: 重試間隔的最大值
					// 無論重試多少次，等待時間都不會超過這個值
					MaxDelay: 30 * time.Second,
				},
			},
		))

	// 創建連接
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	return grpc.DialContext(ctx, config.URL, options...)
}

func InjectClient(params Config) fx.Option {
	return fx.Options(
		fx.Supply(params),
		fx.Provide(fx.Annotate(NewClient, fx.As(new(grpc.ClientConnInterface)))),
	)
}
