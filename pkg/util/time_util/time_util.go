package time_util

import (
	"time"
)

// FormatUnixTime 格式化 Unix 時間戳
// example: FormatUnixTime(1737590400, "2006-01-02 15:04:05")
// example(yyyy-mm-dd): FormatUnixTime(1737590400, "2006-01-02")
func FormatUnixTime(timestamp int64, layout string) string {
	return time.Unix(timestamp/1000, (timestamp%1000)*1e6).Format(layout)
}

// ParseTime 解析時間
// example: ParseTime("2006-01-02 15:04:05", "2021-01-01 00:00:00")
// example(yyyy-mm-dd): ParseTime("2006-01-02", "2021-01-01")
func ParseTime(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// NowString 返回當前時間的字符串
// example: NowString("2006-01-02 15:04:05")
// example(yyyy-mm-dd): NowString("2006-01-02")
func NowString(layout string) string {
	return time.Now().Format(layout)
}

// CurrentEndOfDay 返回當前時間的當天結束時間 (23:59:59.999999999)
func CurrentEndOfDay() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
}
