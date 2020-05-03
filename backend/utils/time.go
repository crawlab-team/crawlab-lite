package utils

import (
	"math"
	"time"
)

func ConvertLocalTime(t time.Time) time.Time {
	return t.In(time.Local)
}

func ConvertTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func ConvertLocalTimeString(t time.Time) string {
	t = ConvertLocalTime(t)
	return ConvertTimeString(t)
}

func ConvertTimestamp(t time.Time) float64 {
	return float64(t.UnixNano()/1e6) / math.Pow10(0)
}

// 秒时间戳
func NowUnix() int64 {
	return time.Now().Unix()
}

// 毫秒时间戳
func NowTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
