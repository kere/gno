package util

import (
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
	return IntZipBaseStr(u64, BaseChars)
}

// IntZipBaseStr int 转换压缩成字符串列表内的字符串
func IntZipBaseStr(num uint64, s []byte) []byte {
	l := uint64(len(s))

	// v, m := calculateZip(l, num)
	v, m := num/l, num%l
	result := make([]byte, 0, int(v)+2)
	result = append(result, s[m])

	for v >= l {
		v, m = v/l, v%l
		// result = append([]byte{s[m]}, result...)
		result = append(result, byte(s[m]))
	}

	if v > 0 {
		// result = append([]byte{s[v]}, result...)
		result = append(result, byte(s[v]))
	}
	return result
}

// // BaseStrToInt10 压缩字符串回退为10位整数
// func BaseStrToInt10(str string, s []byte) (uint64, error) {
// 	l := len(str)
// 	index := bytes.IndexByte(s, str[l-1])
// 	if index < 0 {
// 		return 0, errors.New("parse failed")
// 	}
// 	sum := uint64(index)
// 	base := len(s)
//
// 	for i := 1; i < l; i++ {
// 		index = bytes.IndexByte(s, str[l-i-1])
// 		if index < 0 {
// 			return 0, errors.New("parse failed")
// 		}
// 		sum += uint64(math.Pow(float64(index*base), float64(i)))
// 	}
// 	return sum, nil
// }

// // Unzip62 反向计算压缩字符串
// func Unzip62(str string) (uint64, error) {
// 	return BaseStrToInt10(str, BaseChars)
// }
