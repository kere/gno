package util

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

var (
	// BaseChars 基础字符
	BaseChars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// IntZipTo62 把数字压缩成字符串，基于62字符列表
func IntZipTo62(num uint64) []byte {
	return IntZipBaseStr(num, BaseChars)
}

// IntZipBaseStr int 转换压缩成字符串列表内的字符串
func IntZipBaseStr(num uint64, s []byte) []byte {
	l := uint64(len(s))

	parse := func(n uint64) (uint64, uint64) {
		return n / l, n % l
	}
	result := []byte{}
	v, m := parse(num)
	result = append(result, s[m])

	for v >= l {
		v, m = parse(v)
		result = append([]byte{s[m]}, result...)
	}

	if v > 0 {
		result = append([]byte{s[v]}, result...)
	}
	return result
}

// Round func
func Round(f float64, n int) float64 {
	pow10n := math.Pow10(n)
	if f-float64(int(f)) == 0 {
		return f
	}

	return math.Trunc((f+0.5/pow10n)*pow10n) / pow10n
}

// RandInt 随机整数
func RandInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// ParseIntDefault func
func ParseIntDefault(s string, defaultVal int) int {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultVal
	}

	return int(v)
}

// ParseFloatDefault func
func ParseFloatDefault(s string, defaultVal float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}

	return v
}
