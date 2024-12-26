package convert

import (
	"github.com/spf13/cast"
)

func ToInt(v any) int {
	return cast.ToInt(v)
}
