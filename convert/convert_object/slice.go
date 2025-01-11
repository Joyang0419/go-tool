package convert_object

import (
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// SliceToMap 將 slice 轉換為 map，支援自定義 key 和 value 的轉換
func SliceToMap[T any, K comparable, V any](slice []T, fn func(T) (K, V)) map[K]V {
	return lo.Associate(slice, fn)
}

// SliceToSet 將 slice 轉換為 set
func SliceToSet[T comparable](slice []T) []T {
	return lo.Uniq(slice)
}

func CreateSlice[SpecType any](args ...any) []SpecType {
	slice := make([]SpecType, 0, len(args))

	// 先判斷一次目標類型
	var convertFn func(any) any
	switch any(new(SpecType)).(type) {
	case *int:
		convertFn = func(arg any) any { return cast.ToInt(arg) }
	case *int8:
		convertFn = func(arg any) any { return cast.ToInt8(arg) }
	case *int16:
		convertFn = func(arg any) any { return cast.ToInt16(arg) }
	case *int32:
		convertFn = func(arg any) any { return cast.ToInt32(arg) }
	case *int64:
		convertFn = func(arg any) any { return cast.ToInt64(arg) }
	case *uint:
		convertFn = func(arg any) any { return cast.ToUint(arg) }
	case *uint8:
		convertFn = func(arg any) any { return cast.ToUint8(arg) }
	case *uint16:
		convertFn = func(arg any) any { return cast.ToUint16(arg) }
	case *uint32:
		convertFn = func(arg any) any { return cast.ToUint32(arg) }
	case *uint64:
		convertFn = func(arg any) any { return cast.ToUint64(arg) }
	case *float32:
		convertFn = func(arg any) any { return cast.ToFloat32(arg) }
	case *float64:
		convertFn = func(arg any) any { return cast.ToFloat64(arg) }
	case *string:
		convertFn = func(arg any) any { return cast.ToString(arg) }
	case *bool:
		convertFn = func(arg any) any { return cast.ToBool(arg) }
	case *time.Time:
		convertFn = func(arg any) any { return cast.ToTime(arg) }
	case *time.Duration:
		convertFn = func(arg any) any { return cast.ToDuration(arg) }
	default:
		convertFn = func(arg any) any { return arg }
	}

	// 使用確定的轉換函數處理所有參數
	for _, arg := range args {
		element := convertFn(arg)
		if specTypeElement, ok := element.(SpecType); ok {
			slice = append(slice, specTypeElement)
		}
	}
	return slice
}

func SliceContains[T comparable](slice []T, item T) bool {
	return lo.Contains(slice, item)
}

func SliceToSpecifiedType[T any, SpecType any](slice []T, fn func(T) SpecType) []SpecType {
	return lo.Map(slice, func(item T, _ int) SpecType {
		return fn(item)
	})
}
