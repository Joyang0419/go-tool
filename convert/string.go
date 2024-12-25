package convert

import (
	"github.com/spf13/cast"
)

func StrToInt(s string) int {
	return cast.ToInt(s)
}
