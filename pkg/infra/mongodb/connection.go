package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

type ConnectionConfig struct {
	// MongoDB 連接配置
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`

	// 認證資訊
	Auth struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"auth"`

	/*
		連接池參數建議值：
		- MinPoolSize: 通常設置為預期基礎負載的連接數，例如 5-10
		- MaxPoolSize: 根據服務器資源和負載能力設置，通常 100-200
		- MaxConnIdleTime: 5-10 分鐘 (300s-600s)
		- Timeout: 30 秒，可根據實際需求調整
	*/
	// 連接池配置
	Pool struct {
		// 最小連接池大小
		MinPoolSize uint64 `mapstructure:"min_pool_size"`
		// 最大連接池大小
		MaxPoolSize uint64 `mapstructure:"max_pool_size"`
		// 連接最大空閒時間
		MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
	} `mapstructure:"pool"`

	// 操作超時設置
	Timeout time.Duration `mapstructure:"timeout"`
}

func NewMongoDBConnection(config ConnectionConfig) *mongo.Database {
	ctx := context.Background()

	// 配置客戶端選項
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetAuth(options.Credential{
			Username: config.Auth.Username,
			Password: config.Auth.Password,
		}).
		SetMinPoolSize(config.Pool.MinPoolSize).
		SetMaxPoolSize(config.Pool.MaxPoolSize).
		SetMaxConnIdleTime(config.Pool.MaxConnIdleTime).
		SetTimeout(config.Timeout)

	// 建立連接
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Sprintf("[NewMongoDB]mongo.Connect error: %v", err))
	}

	// 測試連接
	if err = client.Ping(ctx, nil); err != nil {
		panic(fmt.Sprintf("[NewMongoDB]client.Ping error: %v", err))
	}

	return client.Database(config.Database)
}

func InjectMongoDBConnection(config ConnectionConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Provide(NewMongoDBConnection),
	)
}
