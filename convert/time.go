package convert

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// CurrentMonthStartDay 获取当前月的第一天; timezone: example carbon.UTC
// 沒有放的話就是系統預設的時區
func CurrentMonthStartDay(timezone ...string) time.Time {
	return carbon.Now(timezone...).StartOfMonth().StdTime()
}
