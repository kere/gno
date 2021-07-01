package util

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// SplitStrNotSafe 没有内存消耗，返回的内存地址不安全
func SplitStrNotSafe(src string, sep string) []string {
	n := len(sep)
	if len(src) == 0 {
		return nil
	}
	arr := make([]string, 0, 20)
	b := src

	var index int
	for index > -1 {
		index = strings.Index(b, sep)
		if index == -1 {
			break
		}
		if index == 0 {
			arr = append(arr, SEmptyString)
		} else {
			arr = append(arr, b[:index])
		}
		if index+n >= len(b) {
			b = ""
			arr = append(arr, SEmptyString)
			break
		} else {
			b = b[index+n:]
		}
	}

	if len(b) > 0 {
		arr = append(arr, b)
	}
	return arr
}

// SplitBytesNotSafe 没有内存消耗，返回的内存地址不安全
func SplitBytesNotSafe(src []byte, sep []byte) [][]byte {
	n := len(sep)
	if len(src) == 0 {
		return nil
	}
	arr := make([][]byte, 0, 20)
	b := src

	var index int
	for index > -1 {
		index = bytes.Index(b, sep)
		if index == -1 {
			break
		}
		if index == 0 {
			arr = append(arr, nil)
		} else {
			arr = append(arr, b[:index])
		}
		if index+n >= len(b) {
			b = nil
			arr = append(arr, nil)
			break
		} else {
			b = b[index+n:]
		}
	}

	if len(b) > 0 {
		arr = append(arr, b)
	}
	return arr
}

// SplitStr2Floats split string to []float64
func SplitStr2Floats(src, sep string) ([]float64, error) {
	return splitBytes2Floats(Str2Bytes(src), Str2Bytes(sep), false)
}

// SplitBytes2Floats split string to []float64
func SplitBytes2Floats(src, sep []byte) ([]float64, error) {
	return splitBytes2Floats(src, sep, false)
}

// SplitStr2FloatsP split string to []float64
func SplitStr2FloatsP(src, sep string) ([]float64, error) {
	return splitBytes2Floats(Str2Bytes(src), Str2Bytes(sep), true)
}

// SplitBytes2FloatsP split string to []float64
func SplitBytes2FloatsP(src, sep []byte) ([]float64, error) {
	return splitBytes2Floats(src, sep, true)
}

// splitBytes2Floats split string to []float64
func splitBytes2Floats(src, sep []byte, isPool bool) ([]float64, error) {
	arr := SplitBytesNotSafe(src, sep)
	count := len(arr)
	var result []float64
	if isPool {
		result = GetFloats()
	} else {
		result = make([]float64, 0, 20)
	}

	var err error
	var v float64
	for i := 0; i < count; i++ {
		if len(arr[i]) == 0 {
			continue
		}

		switch BytesNumType(arr[i]) {
		case 'f':
			v, err = strconv.ParseFloat(Bytes2Str(arr[i]), 64)
		case 'i':
			var val int64
			val, err = strconv.ParseInt(Bytes2Str(arr[i]), 10, 64)
			v = float64(val)
		default:
			return result, errors.New("do Ints:can not to parse str to num")
		}
		if err != nil {
			return result, err
		}
		result = append(result, v)
	}
	return result, nil
}

// SplitStr2Int64 split string to []int64
func SplitStr2Int64(src, sep string) ([]int64, error) {
	return splitBytes2Int64(Str2Bytes(src), Str2Bytes(sep), false)
}

// SplitBytes2Int64 split string to []int64
func SplitBytes2Int64(src, sep []byte) ([]int64, error) {
	return splitBytes2Int64(src, sep, false)
}

// SplitStr2Int64P split string to []int64
func SplitStr2Int64P(src, sep string) ([]int64, error) {
	return splitBytes2Int64(Str2Bytes(src), Str2Bytes(sep), true)
}

// SplitBytes2Int64P split string to []int64
func SplitBytes2Int64P(src, sep []byte) ([]int64, error) {
	return splitBytes2Int64(src, sep, true)
}

// splitBytes2Int64 split string to []int64
func splitBytes2Int64(src, sep []byte, isPool bool) ([]int64, error) {
	arr := SplitBytesNotSafe(src, sep)
	count := len(arr)
	var result []int64
	if isPool {
		result = GetInt64s()
	} else {
		result = make([]int64, 0, 20)
	}

	var err error
	var v int64
	for i := 0; i < count; i++ {
		if len(arr[i]) == 0 {
			continue
		}

		switch BytesNumType(arr[i]) {
		case 'f':
			var val float64
			val, err = strconv.ParseFloat(Bytes2Str(arr[i]), 64)
			v = int64(val)
		case 'i':
			v, err = strconv.ParseInt(Bytes2Str(arr[i]), 10, 64)
		default:
			return result, errors.New("do Ints:can not to parse str to num")
		}
		if err != nil {
			return result, err
		}
		result = append(result, v)
	}
	return result, nil
}

// SplitStr2Int split string to []int
func SplitStr2Int(src, sep string) ([]int, error) {
	return splitBytes2Int(Str2Bytes(src), Str2Bytes(sep), false)
}

// SplitBytes2Int split string to []int
func SplitBytes2Int(src, sep []byte) ([]int, error) {
	return splitBytes2Int(src, sep, false)
}

// SplitStr2IntP split string to []int
func SplitStr2IntP(src, sep string) ([]int, error) {
	return splitBytes2Int(Str2Bytes(src), Str2Bytes(sep), true)
}

// SplitBytes2IntP split string to []int
func SplitBytes2IntP(src, sep []byte) ([]int, error) {
	return splitBytes2Int(src, sep, true)
}

// splitBytes2Int split string to []int
func splitBytes2Int(src, sep []byte, isPool bool) ([]int, error) {
	arr := SplitBytesNotSafe(src, sep)
	count := len(arr)
	var result []int
	if isPool {
		result = GetInts()
	} else {
		result = make([]int, 0, 20)
	}

	var err error
	var v int
	for i := 0; i < count; i++ {
		if len(arr[i]) == 0 {
			continue
		}

		switch BytesNumType(arr[i]) {
		case 'f':
			var val float64
			val, err = strconv.ParseFloat(Bytes2Str(arr[i]), 64)
			v = int(val)
		case 'i':
			var val int64
			val, err = strconv.ParseInt(Bytes2Str(arr[i]), 10, 64)
			v = int(val)
		default:
			return result, errors.New("do Ints:can not to parse str to num")
		}
		if err != nil {
			return result, err
		}
		result = append(result, v)
	}
	return result, nil
}
