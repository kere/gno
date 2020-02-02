package util

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

var (
	// DBTimeFormat 数据库默认时间格式
	DBTimeFormat = time.RFC3339

	bytesPool sync.Pool
)

// JSONCopy obj1 to obj2
func JSONCopy(from, to interface{}) error {
	src, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(src, to)
}

// EqBytes with 2 []bytes
func EqBytes(arr1, arr2 []byte) bool {
	l := len(arr1)
	if l != len(arr2) {
		return false
	}
	for i := 0; i < l; i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

// // ConvertA2B f
// func ConvertA2B(from, to interface{}) error {
// 	src, err := json.Marshal(from)
// 	if err != nil {
// 		return err
// 	}
// 	return json.Unmarshal(src, to)
// }

var (
	bSeparator = []byte("\\")
	bSlash     = []byte("/")
)

// PathToURLb convert path to url
func PathToURLb(p []byte) []byte {
	return bytes.Replace(p, bSeparator, bSlash, -1)
}

// Str2Bytes to bytes
func Str2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Bytes2Str bytes convert to string
func Bytes2Str(z []byte) string {
	return *(*string)(unsafe.Pointer(&z))
}

// RandStr 生成任意长度的字符串
func RandStr(l int) []byte {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = BaseChars[rand.Intn(62)]
	}
	return bytes
}

// StringsItemCount 在数组中出现的次数
func StringsItemCount(arr []string, val string) int {
	var count int
	for _, v := range arr {
		if v == val {
			count++
		}
	}
	return count
}

// Int64sItemCount 在数组中出现的次数
func Int64sItemCount(arr []int64, val int64) int {
	var count int
	for _, v := range arr {
		if v == val {
			count++
		}
	}
	return count
}

// IsZeroFloats float64s
func IsZeroFloats(arr []float64) bool {
	l := len(arr)
	for i := 0; i < l; i++ {
		if arr[i] != 0 {
			return false
		}
	}
	return true
}
