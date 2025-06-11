package standalone

import (
	"context"
	"fmt"
	"os"
	"time"
)

func New(config Config) *redis.Client {
	host, err := os.Hostname()
	if err != nil {
		msg := fmt.Sprintf("standalone.New os.Hostname error: %v", err)
		panic(msg)
	}

	clientName := host + "-" + config.ClientName
	option := &redis.Options{
		Addr:         config.Address,
		ClientName:   clientName,
		Username:     config.Username,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		// 故意設定成true, 才可以讓使用者用Ctx控制timeout
		ContextTimeoutEnabled: true,
		PoolSize:              config.PoolSize,
		PoolTimeout:           config.PoolTimeout,
		MinIdleConns:          config.MinIdleConns,
		MaxIdleConns:          config.MaxIdleConns,
		ConnMaxIdleTime:       config.ConnMaxIdleTime,
	}

	client := redis.NewClient(option)
	// Ping 測試，5 秒超時
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close() // ping 不到就關閉客戶端
		msg := fmt.Sprintf("standalone.New client.Ping error: %v", err)
		panic(msg)
	}

	if config.Logger != nil {
		redis.SetLogger(config.Logger)
		client.AddHook(config.Logger)
	}

	return client
}
