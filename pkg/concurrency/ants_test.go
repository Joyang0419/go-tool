package concurrency

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithRecovery(t *testing.T) {
	fn := func() {
		time.Sleep(1 * time.Minute)
	}

	WithRecovery(context.TODO(), fn)
	WithRecovery(context.Background(), func() {
		fmt.Println("hello world")
	})

	time.Sleep(3 * time.Second)
}
