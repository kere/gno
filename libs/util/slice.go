package util

import (
	"strconv"
	"strings"
)

// SplitStr2FloatSlice split string to []float64
func SplitStr2FloatSlice(s, sep string) ([]float64, error) {
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
func InStrings(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// InInt64s 是否在数组中出现
func InInt64s(arr []int64, val int64) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// InInts 是否在数组中出现
func InInts(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
