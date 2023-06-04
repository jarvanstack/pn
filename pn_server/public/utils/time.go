package utils

import (
	"fmt"
	"time"
)

const (
	timeFormart = "2006-01-02 15:04:05"
	// 2025-02-19T06:17:49.000Z
	timeFormar2 = "2006-01-02T15:04:05.000Z"
	// 2023-02-19T06:17:49
	timeFormar3 = "2006-01-02_15-04-05"
	zone        = "Asia/Shanghai"
)

func TimeToStr(t time.Time) string {
	return t.Format(timeFormart)
}

func TimeToStr3(t time.Time) string {
	return t.Format(timeFormar3)
}

func StrToTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("时间字符串为空")
	}

	var ret time.Time
	var err error

	ret, err = time.ParseInLocation(timeFormart, s, time.Local)
	if err != nil {
		ret, err = time.ParseInLocation(timeFormar2, s, time.Local)
	}

	return ret, err
}

// 获取一年后的时间
func GetYearAfterTime() time.Time {
	return time.Now().AddDate(1, 0, 0)
}
