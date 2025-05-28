package logger

import (
	"context"
	"sync"
)

/*
	就是調用這邊, 誰實現就用哪個logger, example: slog, zap ...
*/

var (
	// Log 有 Default Log, 透過Init的, 只提供一次性替換
	Log  Interface = DefaultLogger{}
	once sync.Once
)

// Init 一次性初始化 Logger
func Init(logger Interface) {
	onceDoFn := func() {
		Log = logger
	}
	once.Do(onceDoFn)
}

type DefaultLogger struct{}

func (receiver DefaultLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	return
}

func (receiver DefaultLogger) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (receiver DefaultLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return next
}

func (receiver DefaultLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
