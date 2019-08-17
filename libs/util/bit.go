package util

import "strconv"

// GetBitStr for bit setup
func GetBitStr(n int) []byte {
	v := bytesPool.Get()
	var arr []byte
	if v == nil {
		arr = make([]byte, 0, n)
	} else {
		arr = (v.([]byte))[:0]
	}

	for i := 0; i < n; i++ {
		arr = append(arr, '0')
	}
	return arr
}

// PutBytes for bit setup
func PutBytes(arr []byte) {
	bytesPool.Put(arr)
}

// BitStr2Uint uint
func BitStr2Uint(b []byte) uint64 {
	s := BytesToStr(b)
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
