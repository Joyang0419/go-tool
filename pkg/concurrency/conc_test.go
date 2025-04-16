package concurrency

import (
	"testing"
	"time"
)

func TestGoWithRecovery(t *testing.T) {
	fn := func() {
		panic("hello world")
	}
	GoWithRecovery(fn)

	// 等一下log, 確認有收到panic msg
	time.Sleep(1 * time.Second)
}
