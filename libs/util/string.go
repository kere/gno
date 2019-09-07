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
