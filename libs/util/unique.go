package util

import "hash/crc32"

// IID32 return iid
func IID32(str ...string) int64 {
	ieee := crc32.NewIEEE()
	count := len(str)
	for i := 0; i < count; i++ {
		ieee.Write(Str2Bytes(str[i]))
	}
	return int64(ieee.Sum32())
}

// Int64sUnique 获取稀有的整数序列
func Int64sUnique(arr []int64) []int64 {
	m := make(map[int64]int)
	for _, i := range arr {
		m[i] = 1
	}

	u := make([]int64, len(m))
	i := 0
	for k := range m {
		u[i] = k
		i++
	}
	return u
}

// IntsUnique 获取稀有的整数序列
func IntsUnique(arr []int) []int {
	m := make(map[int]int)
	for _, i := range arr {
		m[i] = 1
	}

	u := make([]int, len(m))
	i := 0
	for k := range m {
		u[i] = k
		i++
	}
	return u
}

// StringsUnique 获取稀有的字符串序列
func StringsUnique(arr []string) []string {
	m := make(map[string]int)
	for _, s := range arr {
		m[s] = 1
	}

	u := make([]string, len(m))
	i := 0
	for k := range m {
		u[i] = k
		i++
	}
	return u
}
