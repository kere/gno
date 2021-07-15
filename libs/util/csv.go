package util

import (
	"fmt"
	"io"
	"strconv"
)

func WriteCsvRow(w io.Writer, row []interface{}, decimal int) error {
	count := len(row)
	if count == 0 {
		return nil
	}
	for i := 0; i < count; i++ {
		switch row[i].(type) {
		case string:
			w.Write(Str2Bytes(row[i].(string)))
		case []byte:
			w.Write(row[i].([]byte))
		case nil:
		case int:
			str := strconv.FormatInt(int64(row[i].(int)), 10)
			w.Write(Str2Bytes(str))
		case int32:
			str := strconv.FormatInt(int64(row[i].(int32)), 10)
			w.Write(Str2Bytes(str))
		case int64:
			str := strconv.FormatInt(row[i].(int64), 10)
			w.Write(Str2Bytes(str))
		case float32:
			// str := strconv.FormatFloat(float64(row[i].(float32)), 'f', decimal, 64)
			str := HumanFloat(float64(row[i].(float32)), decimal)
			w.Write(Str2Bytes(str))
		case float64:
			// str := strconv.FormatFloat(row[i].(float64), 'f', decimal, 64)
			str := HumanFloat(row[i].(float64), decimal)
			w.Write(Str2Bytes(str))
		default:
			w.Write(Str2Bytes(fmt.Sprint(row[i])))
		}
		if i < count-1 {
			w.Write(BComma)
		}
	}
	_, err := w.Write(BLineBreak)
	return err
}
