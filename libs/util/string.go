package util

import (
	"bytes"
	"fmt"
)

// CutString 截取字符串
func CutString(str string, length int) string {
	if len(str) <= length {
		return str
	}

	return fmt.Sprint(str[:], "...")
}

// CamelCase name
func CamelCase(name string) string {
	isup := true
	l := len(name)
	if l == 0 {
		return ""
	}

	word := make([]byte, 1)
	src := make([]byte, 0, l)
	for i := 0; i < l; i++ {
		if isup {
			word[0] = name[i]
			src = append(src, bytes.ToUpper(word)...)
			isup = false
			continue
		}
		if name[i] == '_' {
			isup = true
			continue
		}
		src = append(src, name[i])
	}
	return Bytes2Str(src)
}

func IsGBK(data []byte) bool {
	l := len(data)
	var i int = 0
	for i < l {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if i < l-1 && data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				if i < l {
					continue
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}
