package util

import (
	"fmt"
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

// PathToURL convert path to url
func PathToURL(items ...string) string {
	s := filepath.Join(items...)
	if filepath.Separator == '/' {
		return s
	}
	return strings.Replace(s, "\\", "/", -1)
}

// BytesToStr bytes convert to string
func BytesToStr(b []byte) string {
	// return *(*string)(unsafe.Pointer(&s))
	var s string
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return s
}

// StrToBytes bytes convert to string
func StrToBytes(s string) []byte {
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

// CutString 截取字符串
func CutString(str string, length int) string {
	if len(str) <= length {
		return str
	}

	return fmt.Sprint(str[:], "...")
}
