package logger

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Interface interface {
	// Printf 給 redis.SetLogger(config.Logger);
	/*
		    可以看到這些log
			ERROR redis internal: connection lost
			WARN redis internal: connection timeout
	*/
	Printf(ctx context.Context, format string, v ...interface{})

	/*
		redis.Hook; 可以看到每個命令執行
	*/
	redis.Hook
}
