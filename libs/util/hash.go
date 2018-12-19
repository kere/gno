package util

import (
	"crypto/md5"
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

//MD5 md5校验
func MD5(src []byte) []byte {
	m := md5.New()
	m.Write(src)
	return m.Sum(nil)
}

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
	return IntZipBaseStr(u64, BaseChars)
}

func calculateZip(l, n uint64) (uint64, uint64) {
	return n / l, n % l
}

// IntZipBaseStr int 转换压缩成字符串列表内的字符串
func IntZipBaseStr(num uint64, s []byte) []byte {
	l := uint64(len(s))

	result := []byte{}
	v, m := calculateZip(l, num)
	result = append(result, s[m])

	for v >= l {
		v, m = calculateZip(l, v)
		result = append([]byte{s[m]}, result...)
	}

	if v > 0 {
		result = append([]byte{s[v]}, result...)
	}
	return result
}
