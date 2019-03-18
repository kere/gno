package util

import (
	"sort"
	"strconv"
	"strings"
)

// CopyStrings copy
func CopyStrings(arr []string) []string {
	l := len(arr)
	if l == 0 {
		return []string{}
	}
	src := make([]string, l)
	copy(arr, src)
	return src
}

// CopyInts copy
func CopyInts(arr []int) []int {
	l := len(arr)
	if l == 0 {
		return []int{}
	}
	src := make([]int, l)
	copy(arr, src)
	return src
}

// SplitStr2Floats split string to []float64
func SplitStr2Floats(s, sep string) ([]float64, error) {
	arr := strings.Split(s, sep)
	var err error
	var v float64
	l := len(arr)
	result := make([]float64, 0)

	for i := 0; i < l; i++ {
		v, err = strconv.ParseFloat(arr[i], 64)
		if err != nil {
			continue
		}
		result = append(result, v)
	}
	return result, nil
}

// InStrings 是否在数组中出现
func InStrings(val string, arr []string) bool {
	l := len(arr)
	if l < 500 {
		for _, v := range arr {
			if v == val {
				return true
			}
		}
		return false
	}

	tmp := make([]string, l)
	copy(tmp, arr)
	sort.Strings(tmp)
	index := sort.SearchStrings(tmp, val)
	return index != l
}

// // InInt64s 是否在数组中出现
// func InInt64s(val int64, arr []int64) bool {
// 	for _, v := range arr {
// 		if v == val {
// 			return true
// 		}
// 	}
// 	return false
// }

// InInts 是否在数组中出现
func InInts(val int, arr []int) bool {
	l := len(arr)
	if l < 500 {
		for _, v := range arr {
			if v == val {
				return true
			}
		}
		return false
	}
	tmp := make([]int, l)
	copy(tmp, arr)
	sort.Ints(tmp)
	index := sort.SearchInts(tmp, val)
	return index != l
}

// InFloats 是否在数组中出现
func InFloats(val float64, arr []float64) bool {
	l := len(arr)
	if l < 500 {
		for _, v := range arr {
			if v == val {
				return true
			}
		}
		return false
	}
	tmp := make([]float64, l)
	copy(tmp, arr)
	sort.Float64s(tmp)
	index := sort.SearchFloat64s(tmp, val)
	return index != l
}

// SameStrings 是否数组相同
func SameStrings(arr1, arr2 []string) bool {
	l := len(arr1)
	if l != len(arr2) {
		return false
	}
	tmp1 := make([]string, l)
	tmp2 := make([]string, l)
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
	tmp1 := make([]int, l)
	tmp2 := make([]int, l)
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
