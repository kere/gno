package util

import (
	"hash/crc32"
	"sort"
)

// IID32 return iid
func IID32(str ...string) int64 {
	ieee := crc32.NewIEEE()
	count := len(str)
	for i := 0; i < count; i++ {
		ieee.Write(Str2Bytes(str[i]))
	}
	return int64(ieee.Sum32())
}

// Int64sUniqueP 获取稀有的整数序列
func Int64sUniqueP(arr []int64) []int64 {
	l := len(arr)
	if l == 0 {
		return arr
	}
	tmp := Int64sOrder(GetInt64s(l))
	defer PutInt64s(tmp)
	copy(tmp, arr)
	tmp.Sort()

	// u := make([]int64, len(m))
	u := GetInt64s(1, l)
	v := tmp[0]
	u[0] = v

	for i := 1; i < l; i++ {
		if v == tmp[i] {
			continue
		}
		v = tmp[i]
		u = append(u, v)
	}
	return u
}

// IntsUniqueP 获取稀有的整数序列
func IntsUniqueP(arr []int) []int {
	l := len(arr)
	if l == 0 {
		return arr
	}
	tmp := GetInts(l)
	defer PutInts(tmp)

	copy(tmp, arr)
	sort.Ints(tmp)

	// u := make([]int64, len(m))
	u := GetInts(1, l)
	v := tmp[0]
	u[0] = v

	for i := 1; i < l; i++ {
		if v == tmp[i] {
			continue
		}
		v = tmp[i]
		u = append(u, v)
	}
	return u
}

// StringsUniqueP 获取稀有的字符串序列
func StringsUniqueP(arr []string) []string {
	l := len(arr)
	if l == 0 {
		return arr
	}
	tmp := GetStrings(l)
	defer PutStrings(tmp)

	copy(tmp, arr)
	sort.Strings(tmp)

	// u := make([]int64, len(m))
	u := GetStrings(1, l)
	v := tmp[0]
	u[0] = v

	for i := 1; i < l; i++ {
		if v == tmp[i] {
			continue
		}
		v = tmp[i]
		u = append(u, v)
	}
	return u
}
