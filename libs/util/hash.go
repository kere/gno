package util

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"time"

	"github.com/spaolacci/murmur3"
)

var (
	// BaseChars 基础字符
	BaseChars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// ECMATable table
var ECMATable = crc64.MakeTable(crc64.ECMA)

//CRC64Token crc校验
func CRC64Token(src []byte) []byte {
	csum := crc64.Checksum(src, ECMATable)
	return IntZipTo62(csum)
}

//CRC32Token crc校验
func CRC32Token(src []byte) []byte {
	ieee := crc32.NewIEEE()
	// io.WriteString(ieee, str)
	ieee.Write(src)
	v64 := uint64(ieee.Sum32())
	return IntZipTo62(v64)
}

// Unique 获得一个稀有字符
func Unique() string {
	u64 := time.Now().UTC().UnixNano()
	return UUIDshort(u64)
}

// UUIDshort 获得一个值的Short UUID
func UUIDshort(v interface{}) string {
	v64 := murmur3.Sum64([]byte(fmt.Sprint(v)))
	return string(IntZipTo62(v64))
}

// IntZipTo62 把数字压缩成字符串，基于62字符列表
func IntZipTo62(u64 uint64) []byte {
	return IntZipTo(u64, BaseChars)
}

// IntZipTo int 转换压缩成字符串列表内的字符串
func IntZipTo(num uint64, table []byte) []byte {
	l := uint64(len(table))
	if num < l {
		return []byte{table[num]}
	}

	// v, m := calculateZip(l, num)
	v, m := num/l, num%l
	result := make([]byte, 0, 12)
	result = append(result, table[m])

	for v >= l {
		v, m = v/l, v%l
		result = append(result, byte(table[m]))
	}

	if v > 0 {
		result = append(result, byte(table[v]))
	}
	return result
}

// UnIntZip int
func UnIntZip(s string, table []byte) int64 {
	l := len(s)
	if l == 0 {
		return -1
	}
	if l == 1 {
		k := bytes.IndexRune(table, rune(s[0]))
		return int64(k)
	}

	n := int64(len(table))
	k := bytes.IndexRune(table, rune(s[l-1]))
	if k < 0 {
		return -1
	}
	val := int64(k) * n
	k = bytes.IndexRune(table, rune(s[l-2]))
	if k < 0 {
		return -1
	}
	val += int64(k)

	for i := l - 3; i > -1; i-- {
		k = bytes.IndexRune(table, rune(s[i]))
		if k < 0 {
			return -1
		}
		val = n*val + int64(k)
	}

	return val
}
