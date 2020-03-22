package util

import (
	"strconv"
)

// BitStr2Uint uint
func BitStr2Uint(b []byte) uint64 {
	l := len(b)
	for i := 0; i < l; i++ {
		if b[i] == 0 {
			b[i] = '0'
		}
	}
	s := Bytes2Str(b)
	v, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return 0
	}
	return v
}

// SetBitStrTrue []byte
func SetBitStrTrue(b []byte, i int) {
	l := len(b)
	if i >= l {
		return
	}
	// 反方向排列
	b[l-i-1] = '1'
}

// Float2BitStr string
func Float2BitStr(v float64) string {
	return strconv.FormatUint(uint64(v), 2)
}

// IsTrueAtBitUint uint
func IsTrueAtBitUint(u uint64, i int) bool {
	s := strconv.FormatUint(u, 2)
	l := len(s)
	if i >= l {
		return false
	}
	// 反方向排列
	return s[l-i-1] == '1'
}
