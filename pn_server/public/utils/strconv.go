package utils

import (
	"fmt"
	"strconv"
)

func StrToUint(s string) (uint, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	return uint(num), err
}

func UintToStr(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

func StrToInt(s string) (int, error) {
	num, err := strconv.ParseInt(s, 10, 64)
	return int(num), err
}

func IntToStr(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// 字节数转换为可读字符串
// 比如 1024 转换为 1KB
// 1024 * 1024 + 1024 转换为 1MB
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
