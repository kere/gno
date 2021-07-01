package util

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// StrNumType return i,f,s,
func StrNumType(s string) rune {
	b := Str2Bytes(s)
	return BytesNumType(b)
}

// BytesNumType return i,f,s,
func BytesNumType(b []byte) rune {
	l := len(b)
	var skip0 bool
	// 1: 首字母不是数字
	if b[0] < 46 || 57 < b[0] {
		if l == 1 {
			return 's'
		}
		lowS := strings.ToLower(Bytes2Str(b))
		if lowS == "inf" || lowS == "+inf" || lowS == "-inf" || lowS == "nan" {
			return 'f'
		}
		// 是否首字母是：+ -
		if b[0] != 43 && b[0] != 45 {
			// 首字母是+ -时，1号位置是否为字符
			if b[1] < 46 || 57 < b[1] {
				return 's'
			}
		}
		// 跳过首字符+-
		skip0 = true
	}

	count := len(b)
	start := 0
	if skip0 {
		// 跳过首字符+-
		start = 1
	}
	dot, e, plus := -1, -1, -1
	for i := start; i < count; i++ {
		if b[i] < 47 || 57 < b[i] {
			switch b[i] {
			case 46: // dot .
				if dot != -1 {
					// dot出现2次
					return 's'
				}
				dot = i
			case 69, 101: // E, e
				if e > 0 {
					// dot出现2次
					return 's'
				}
				e = i
			case 43, 45:
				if plus > 0 {
					// dot出现2次
					return 's'
				}
				plus = i
			default:
				return 's'
			}
		}
	}
	if e == -1 && plus == -1 {
		if dot == -1 {
			return 'i'
		} else {
			return 'f'
		}
	}
	if e > 0 && plus == e+1 && l > plus+1 {
		return 'f'
	}

	return 's'
}

// Max int
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// Min int
func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// Max64 int64
func Max64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

// Min64 int64
func Min64(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

// AbsInt64 int64
func AbsInt64(val int64) int64 {
	if val < 0 {
		return -1 * val
	}
	return val
}

// Abs int
func Abs(val int) int {
	if val < 0 {
		return -1 * val
	}
	return val
}

// AbsFloat32 int
func AbsFloat32(val float32) float32 {
	if val < 0 {
		return -1 * val
	}
	return val
}

// Money string
func Money(val float64) string {
	val = Round(val, 2)
	return HumanFloatC(val)
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

// Round32 func
func Round32(val float32, n int) float32 {
	v := 1.0
	if n > 0 {
		v = math.Pow10(n)
	}
	if val < 0 {
		return float32(math.Ceil(float64(val)*v-0.5) / v)
	}

	return float32(math.Floor(float64(val)*v+0.5) / v)
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

// IsEmptyFloats float64s
func IsEmptyFloats(arr []float64) bool {
	l := len(arr)
	for i := 0; i < l; i++ {
		if arr[i] != 0 {
			return false
		}
	}
	return true
}

// HasNaNInFloats float64s
func HasNaNInFloats(arr []float64) bool {
	l := len(arr)
	for i := 0; i < l; i++ {
		if math.IsNaN(arr[i]) {
			return true
		}
	}
	return false
}

// HumanFloatC 用中国4进位方式，显示数字
func HumanFloatC(v float64, args ...int) string {
	prec := -1
	if len(args) != 0 {
		prec = args[0]
	}
	var str string
	str = strconv.FormatFloat(v, 'f', prec, 64)

	index := strings.IndexByte(str, '.')
	if index == -1 {
		b := recopStepN(Str2Bytes(str), 4, ',')
		return Bytes2Str(b)
	}
	// 整数部分
	src := recopStepN(Str2Bytes(str[:index]), 4, ',')
	count := len(str)
	for i := index; i < count; i++ {
		src = append(src, str[i])
	}
	return Bytes2Str(src)
}

// recopStepN 用中国4进位方式，显示数字
func recopStepN(src []byte, stepN int, sep byte) []byte {
	k := 0
	bf := GetBytes()
	count := len(src)
	for i := count - 1; i > -1; i-- {
		if k == stepN {
			k = 0
			bf = append(bf, sep)
		}
		bf = append(bf, src[i])
		k++
	}
	arr := make([]byte, 0, 30)

	count = len(bf)
	for i := count - 1; i > -1; i-- {
		arr = append(arr, bf[i])
	}
	PutBytes(bf)
	return arr
}

// HumanInt64C 用中国4进位方式，显示数字
func HumanInt64C(v int64) string {
	str := fmt.Sprint(v)
	b := recopStepN(Str2Bytes(str), 4, ',')
	return Bytes2Str(b)
}
