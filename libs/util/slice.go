package util

import (
	"sort"
)

// StringsI 在数组中出现的index
func StringsI(val string, arr []string) int {
	l := len(arr)
	for i := 0; i < l; i++ {
		if arr[i] == val {
			return i
		}
	}
	return -1
}

// InStrings 是否在数组中出现
func InStrings(val string, arr []string) bool {
	l := len(arr)
	if l < 300 {
		for _, v := range arr {
			if v == val {
				return true
			}
		}
		return false
	}

	// tmp := make([]string, l)
	tmp := GetStrings(l)
	copy(tmp, arr)
	sort.Strings(tmp)
	index := sort.SearchStrings(tmp, val)
	PutStrings(tmp)
	return index != l
}

// InInts 是否在数组中出现
func InInts(val int, arr []int) bool {
	l := len(arr)
	if l < 300 {
		for _, v := range arr {
			if v == val {
				return true
			}
		}
		return false
	}
	// tmp := make([]int, l)
	tmp := GetInts(l)
	copy(tmp, arr)
	sort.Ints(tmp)
	index := sort.SearchInts(tmp, val)
	PutInts(tmp)
	return index != l
}

// InFloats 是否在数组中出现
func InFloats(val float64, arr []float64) bool {
	l := len(arr)
	if l < 300 {
		for _, v := range arr {
			if v == val {
				return true
			}
		}
		return false
	}
	// tmp := make([]float64, l)
	tmp := GetFloats(l)
	copy(tmp, arr)
	sort.Float64s(tmp)
	index := sort.SearchFloat64s(tmp, val)
	PutFloats(tmp)
	return index != l
}

// SameStrings 是否数组相同
func SameStrings(arr1, arr2 []string) bool {
	l := len(arr1)
	if l != len(arr2) {
		return false
	}
	// tmp1 := make([]string, l)
	// tmp2 := make([]string, l)
	tmp1 := GetStrings(l)
	defer PutStrings(tmp1)
	tmp2 := GetStrings(l)
	defer PutStrings(tmp2)
	copy(tmp1, arr1)
	copy(tmp2, arr2)
	sort.Strings(tmp1)
	sort.Strings(tmp2)

	for i := 0; i < l; i++ {
		if tmp1[i] != tmp2[i] {
			return false
		}
	}

	return true
}

// SameInts 是否数组相同
func SameInts(arr1, arr2 []int) bool {
	l := len(arr1)
	if l != len(arr2) {
		return false
	}
	// tmp1 := make([]int, l)
	// tmp2 := make([]int, l)
	tmp1 := GetInts(l)
	defer PutInts(tmp1)
	tmp2 := GetInts(l)
	defer PutInts(tmp2)
	copy(tmp1, arr1)
	copy(tmp2, arr2)
	sort.Ints(tmp1)
	sort.Ints(tmp2)

	for i := 0; i < l; i++ {
		if tmp1[i] != tmp2[i] {
			return false
		}
	}

	return true
}

// RangeFloats 取一段
func RangeFloats(arr []float64, a, b int) []float64 {
	l := len(arr)
	if a < 0 {
		a = 0
	}
	if b == 0 || b > l-1 {
		b = l - 1
	}
	return arr[a : b+1]
}

// RangeInts 取一段
func RangeInts(arr []int, a, b int) []int {
	l := len(arr)
	if a < 0 {
		a = 0
	}
	if b == 0 || b > l-1 {
		b = l - 1
	}
	return arr[a : b+1]
}

// EachPage
func EachPage(count, pageSize int, f func(pageN, a, b int) bool) int {
	var b, pageN int
	for i := 0; i < count; i++ {
		b = i + pageSize
		if b > count {
			b = count
		}
		pageN++
		if !f(pageN, i, b) {
			return -1
		}
		if b == count {
			return pageN
		}
		i = b - 1
	}

	return -1
}
