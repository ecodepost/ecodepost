package utils

import (
	"time"
)

func GetNowMinTimestamps() int64 {
	local, _ := time.LoadLocation("Local")
	str := time.Now().Format("20060102 15:04:00")
	tt, _ := time.ParseInLocation("20060102 15:04:00", str, local)
	ts := tt.Unix()
	return ts
}

func GetNowDayTimestamps() int64 {
	local, _ := time.LoadLocation("Local")
	str := time.Now().Format("20060102")
	tt, _ := time.ParseInLocation("20060102", str, local)
	ts := tt.Unix()
	return ts
}

// FormatDate 获取当前的日期
func FormatDate(now time.Time) string {
	return now.Format("2006-01-02")
}
