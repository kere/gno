package util

import (
	"strconv"
)

// MaskBytes2Int int
func MaskBytes2Int(b []byte) int {
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
	return int(v)
}

// SetBytesMaskTrue []byte
func SetBytesMaskTrue(b []byte, i int) {
	l := len(b)
	if i >= l {
		return
	}
	// 反方向排列
	b[l-i-1] = '1'
}

// SetIntMask set int mask
func SetIntMask(v, i int, isTrue bool) int {
	s := strconv.FormatInt(int64(v), 2)
	l := len(s)
	if i >= l {
		return 0
	}
	b := Str2Bytes(s)
	if isTrue {
		b[l-1-i] = '1'
	} else {
		b[l-1-i] = '0'
	}
	v2, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return 0
	}
	return int(v2)
}

// IsMaskTrueAt uint 从右向左
func IsMaskTrueAt(u int, i int) bool {
	s := strconv.FormatInt(int64(u), 2)
	l := len(s)
	if i >= l {
		return false
	}
	// 反方向排列
	return s[l-1-i] == '1'
}
