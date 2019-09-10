package util

import (
	"math/rand"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var (
	// DBTimeFormat 数据库默认时间格式
	DBTimeFormat = time.RFC3339

	bytesPool sync.Pool
)

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

// PathToURL convert path to url
func PathToURL(items ...string) string {
	s := filepath.Join(items...)
	if filepath.Separator == '/' {
		return s
	}
	return strings.Replace(s, "\\", "/", -1)
}

// Bytes2Str bytes convert to string
func Bytes2Str(b []byte) string {
	// return *(*string)(unsafe.Pointer(&s))
	var s string
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return s
}

// Str2Bytes bytes convert to string
func Str2Bytes(s string) []byte {
	var b []byte
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len
	return b
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
