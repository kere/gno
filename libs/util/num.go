package util

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// Money string
func Money(val float64) string {
	val = Round(val, 2)
	return humanize.Commaf(val)
}

// Round func
func Round(val float64, n int) float64 {
	v := 1.0
	if n > 0 {
		v = math.Pow10(n)
	}
	if val < 0 {
		return math.Ceil(val*v-0.5) / v
	}

	return math.Floor(val*v+0.5) / v
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
