package convert_number

import (
	"github.com/spf13/cast"
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func ToNumber[T Number](i interface{}) T {
	var result T
	switch any(result).(type) {
	case int:
		return T(cast.ToInt(i))
	case int64:
		return T(cast.ToInt64(i))
	case int32:
		return T(cast.ToInt32(i))
	case float64:
		return T(cast.ToFloat64(i))
	case float32:
		return T(cast.ToFloat32(i))
	case uint:
		return T(cast.ToUint(i))
	case uint64:
		return T(cast.ToUint64(i))
	default:
		return result
	}
}
