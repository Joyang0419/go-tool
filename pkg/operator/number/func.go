package number

import (
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

type TNumber interface {
	constraints.Integer | constraints.Float
}

// InRange 檢查一個數字是否落在指定的範圍內（包括邊界）。
//
// 參數:
//   - current: 要檢查的數字。
//   - from: 範圍的最小值。
//   - to: 範圍的最大值。
//
// 回傳值:
//   - bool: 如果數字落在範圍內（包括邊界），則回傳 true；否則回傳 false。
//
// 範例:
//
//	InRange(5, 1, 10) -> true (5 在範圍 [1, 10] 內)
//	InRange(15, 1, 10) -> false (15 不在範圍 [1, 10] 內)
//	InRange(1, 1, 10) -> true (1 等於範圍最小值)
//	InRange(10, 1, 10) -> true (10 等於範圍最大值)
func InRange[T TNumber](current, from, to T) bool {
	return current >= from && current <= to
}

// Max 回傳給定切片中的最大值。
//
// 此泛型方法適用於任何實現了 TNumber 的類型，
// 如整數、浮點數和字串等具有比較能力的類型。
//
// 參數:
//
//	values - 一個具有 constraints.Ordered 的泛型切片，用來比較並找出最大值。
//
// 回傳值:
//
//	返回切片中具有最大值的元素。如果切片為空，將回傳該類型的零值。
//
// 範例:
//
//	找出整數中的最大值:
//	  maxInt := Max([]int{1, 23, 42, 4}) // maxInt = 42
//
//	找出浮點數中的最大值:
//	  maxFloat := Max([]float64{3.2, 1.5, 4.7, 2.1}) // maxFloat = 4.7
func Max[T TNumber](values ...T) T {
	return lo.Max(values)
}

// Min 回傳給定切片中的最小值。
//
// 此泛型方法適用於任何實現了 constraints.Ordered 的類型，
// 如整數、浮點數和字串等具有比較能力的類型。
//
// 參數:
//
//	values - 一個具有 constraints.Ordered 的泛型切片，用來比較並找出最小值。
//
// 回傳值:
//
//	返回切片中具有最小值的元素。如果切片為空，將回傳該類型的零值。
//
// 範例:
//
//	找出整數中的最小值:
//	  minInt := Min([]int{1, 23, 42, 4}) // minInt = 1
//
//	找出浮點數中的最小值:
//	  minFloat := Min([]float64{3.2, 1.5, 4.7, 2.1}) // minFloat = 1.5
func Min[T TNumber](values ...T) T {
	return lo.Min(values)
}

func Average[T TNumber](values ...T) T {
	return lo.Mean(values)
}

func Sum[T TNumber](values ...T) T {
	return lo.Sum(values)
}
