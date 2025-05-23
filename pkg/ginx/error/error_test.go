package error

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestName(t *testing.T) {
	err1 := ErrBadParam
	err2 := ErrBadParam

	// 確認下這樣不會影響到err2
	fmt.Println(err1.SetTraceID("hello world"))
	fmt.Println(err2)
}

func TestWrap(t *testing.T) {
	err2 := ErrBadParam.SetActualError(Err())
	err3 := err2.SetActualError(Err1())
	fmt.Println(err2.Error())
	fmt.Println(err3.Error())
}

func Err() error {
	return errors.New("err0")
}

func Err1() error {
	return errors.New("err1")
}
