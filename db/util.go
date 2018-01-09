package db

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"strings"
)

func MD5(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func InStrings(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func Arr2InCondition(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	tmp := make([]string, len(arr))
	for i, c := range arr {
		tmp[i] = fmt.Sprint("'", c, "'")
	}
	return strings.Join(tmp, ",")
}

func arrayBaseType(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		return arrayBaseType(typ.Elem())

	default:
		return typ
	}
}
