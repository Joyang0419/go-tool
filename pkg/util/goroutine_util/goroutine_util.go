package goroutine_util

import (
	"log/slog"
)

func GoWithRecover(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("[GoWithRecover]Recovered", slog.Any("error", r))
			}
		}()
		f()
	}()
}
