package logger

import (
	"log/slog"
	"os"
	"sync"
)

/*
	就是調用這邊, 誰實現就用哪個logger, example: slog, zap ...
*/

var (
	// Log 有 Default Log, 透過Init的, 一次性替換
	Log Interface = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}),
	)
	once sync.Once
)

// Init 一次性初始化 Logger
func Init(logger Interface) {
	onceDoFn := func() {
		Log = logger
	}
	once.Do(onceDoFn)
}
