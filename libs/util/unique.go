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

// Int64sUnique 获取稀有的整数序列
func Int64sUnique(arr []int64, isPool bool) []int64 {
	l := len(arr)
	if l == 0 {
		return arr
	}
	var tmp []int64
	if isPool {
		tmp = GetInt64s(l)
		defer PutInt64s(tmp)
	} else {
		tmp = make([]int64, l)
	}
	copy(tmp, arr)
	sorted := Int64sOrder(tmp)
	sorted.Sort()

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

// IntsUnique 获取稀有的整数序列
func IntsUnique(arr []int, isPool bool) []int {
	l := len(arr)
	if l == 0 {
		return arr
	}
	var tmp []int
	if isPool {
		tmp = GetInts(l)
		defer PutInts(tmp)
	} else {
		tmp = make([]int, l)
	}

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

// StringsUnique 获取稀有的字符串序列
func StringsUnique(arr []string, isPool bool) []string {
	l := len(arr)
	if l == 0 {
		return arr
	}
	var tmp []string
	if isPool {
		tmp = GetStrings(l)
		defer PutStrings(tmp)
	} else {
		tmp = make([]string, l)
	}

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

// StringsMergUniq 合并多个数字，获取稀有的字符串序列
func StringsMergUniq(isPool bool, arr ...[]string) []string {
	count := len(arr)
	if count == 0 {
		return nil
	}
	tmp := GetStrings(len(arr[0]))
	defer PutStrings(tmp)
	copy(tmp, arr[0])

	for i := 1; i < count; i++ {
		tmp = append(tmp, arr[i]...)
	}

	return StringsUnique(tmp, isPool)
}
