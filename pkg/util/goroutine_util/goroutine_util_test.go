package goroutine_util

import (
	"testing"
)

func TestGoWithRecover(t *testing.T) {
	fn := func() {
		panic("test panic")
	}

	GoWithRecover(fn)
}
