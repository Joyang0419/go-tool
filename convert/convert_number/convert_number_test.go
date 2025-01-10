package convert_number

import (
	"fmt"
	"testing"
)

func TestToNumber(t *testing.T) {
	fmt.Println(ToNumber[int](true))
	fmt.Println(ToNumber[int64]("33") + 1)
	a := ToNumber[float64]("1.2")
	fmt.Println(a + 2.2)
	fmt.Println(ToNumber[float32](1.2) + 2.2)
}
