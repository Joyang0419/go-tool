package operator

import (
	"golang.org/x/exp/constraints"
)

// NumberInRange 檢查數字是否在範圍內
func NumberInRange[T constraints.Ordered](current, from, to T) bool {
	return current >= from && current <= to
}
