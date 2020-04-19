package util

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

const (
	// PrintLine for print
	PrintLine = "-----------------------"
	// PrintDotted for print
	PrintDotted = "......................."
	// PrintEquals for print
	PrintEquals = "======================="
)

func PrintObj(obj interface{}) {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	n := typ.NumField()
	for i := 0; i < n; i++ {
		ftyp := typ.Field(i)
		name := ftyp.Tag.Get("name")
		if name == "" {
			name = ftyp.Name
		}
		hasPer := strings.HasSuffix(name, "%")

		str := ""
		fval := val.Field(i)
		switch ftyp.Type.Kind() {
		case reflect.Float64:
			v := fval.Float()
			if math.IsInf(v, 0) {
				str = "~"
			} else if math.IsNaN(v) {
				str = "-"
			} else if v == 0 {
				str = "0"
			} else {
				if hasPer {
					str = humanize.Ftoa(v*100) + "%"
				} else {
					str = humanize.Ftoa(v)
				}
			}
		case reflect.Int, reflect.Int64:
			if hasPer {
				str = fmt.Sprint(fval.Int()*100) + "%"
			} else {
				str = fmt.Sprint(fval.Int())
			}
		case reflect.String:
			str = fval.String()
		case reflect.Slice:
			v := fval.Interface()
			switch fval.Interface().(type) {
			case []byte:
				str = Bytes2Str(v.([]byte))
			default:
				str = fmt.Sprint(v)
			}

		default:
			str = fmt.Sprintf("%v", fval.Interface())
		}

		switch ftyp.Tag.Get("pstyle") {
		case "inline":
			fmt.Print(name, ": ", str, "\t\t")
		default:
			fmt.Printf("%s: %s\n", name, str)
		}
		switch ftyp.Tag.Get("print") {
		case "break":
			fmt.Println()
		case "break2":
			fmt.Print("\n\n")
		case "line":
			fmt.Println(PrintLine)
		case "dotted":
			fmt.Println(PrintDotted)
		case "equal":
			fmt.Println(PrintEquals)
		}
	}
	fmt.Println()
}
