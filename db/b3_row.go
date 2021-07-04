package db

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/kere/gno/libs/util"
)

// DBRow
type DBRow struct {
	Values []interface{}
	Fields []string
}

// IsEmpty
func (d *DBRow) IsEmpty() bool {
	return len(d.Values) == 0
}

// Int
func (d *DBRow) Int(field string) int {
	v := d.Int64(field)
	return int(v)
}

// IntAt
func (d *DBRow) IntAt(i int) int {
	v := d.Int64At(i)
	return int(v)
}

// Int64At
func (d *DBRow) Int64(field string) int64 {
	i := util.StringsI(field, d.Fields)
	if i < 0 {
		return 0
	}
	return d.Int64At(i)
}

// Int64
func (d *DBRow) Int64At(i int) int64 {
	typ := reflect.TypeOf(d.Values[i])
	val := reflect.ValueOf(d.Values[i])
	switch typ.Kind() {
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return val.Int()
	case reflect.Slice:
		if typ.String() == "[]uint8" {
			src := []byte(d.Values[i].([]uint8))
			switch util.BytesNumType(src) {
			case 'f':
				v, _ := strconv.ParseFloat(util.Bytes2Str(src), 64)
				return int64(v)
			case 'i':
				v, _ := strconv.ParseInt(util.Bytes2Str(src), 10, 64)
				return v
			}
		}

		return 0
	case reflect.Float64, reflect.Float32:
		return int64(val.Float())
	case reflect.String:
		switch util.BytesNumType(util.Str2Bytes(val.String())) {
		case 'f':
			v, _ := strconv.ParseFloat(val.String(), 64)
			return int64(v)
		case 'i':
			v, _ := strconv.ParseInt(val.String(), 10, 64)
			return v
		}
	}
	return 0
}

// Float32
func (d *DBRow) Float32(field string) float32 {
	v := d.Float(field)
	return float32(v)
}

// Float64
func (d *DBRow) Float(field string) float64 {
	i := util.StringsI(field, d.Fields)
	if i < 0 {
		return 0
	}
	return d.FloatAt(i)
}

// Float64At
func (d *DBRow) FloatAt(i int) float64 {
	typ := reflect.TypeOf(d.Values[i])
	val := reflect.ValueOf(d.Values[i])
	switch typ.Kind() {
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return float64(val.Int())
	case reflect.Slice:
		if typ.String() == "[]uint8" {
			src := []byte(d.Values[i].([]uint8))
			v, _ := strconv.ParseFloat(util.Bytes2Str(src), 64)
			return v
		}
		return 0
	case reflect.Float64, reflect.Float32:
		return val.Float()
	case reflect.String:
		v, _ := strconv.ParseFloat(val.String(), 64)
		return v
	}
	return 0
}

// String
func (d *DBRow) String(field string) string {
	i := util.StringsI(field, d.Fields)
	if i < 0 {
		return ""
	}
	return d.StringAt(i)
}

// StringAt
func (d *DBRow) StringAt(i int) string {
	switch d.Values[i].(type) {
	case []byte:
		return util.Bytes2Str(d.Values[i].([]byte))
	case string:
		return d.Values[i].(string)
	}
	return fmt.Sprint(d.Values[i])
}

// Bytes
func (d *DBRow) Bytes(field string) []byte {
	i := util.StringsI(field, d.Fields)
	if i < 0 {
		return nil
	}
	return d.BytesAt(i)
}

// BytesAt
func (d *DBRow) BytesAt(i int) []byte {
	switch d.Values[i].(type) {
	case []byte:
		return d.Values[i].([]byte)
	case string:
		return util.Str2Bytes(d.Values[i].(string))
	}
	return util.Str2Bytes(fmt.Sprint(d.Values[i]))
}

func (d *DBRow) Int64sAt(i int) ([]int64, error) {
	return Current().Driver.Int64s(d.Values[i].([]byte))
}
func (d *DBRow) Int64sAtP(i int) ([]int64, error) {
	return Current().Driver.Int64sP(d.Values[i].([]byte))
}

func (d *DBRow) StringsAt(i int) ([]string, error) {
	return Current().Driver.Strings(d.Values[i].([]byte))
}
func (d *DBRow) StringsNotSafeAt(i int) ([]string, error) {
	return Current().Driver.StringsNotSafe(d.Values[i].([]byte))
}

func (d *DBRow) IntsAt(i int) ([]int, error) {
	return Current().Driver.Ints(d.Values[i].([]byte))
}
func (d *DBRow) IntsAtP(i int) ([]int, error) {
	return Current().Driver.IntsP(d.Values[i].([]byte))
}

func (d *DBRow) FloatsAt(i int) ([]float64, error) {
	return Current().Driver.Floats(d.Values[i].([]byte))
}
func (d *DBRow) FloatsAtP(i int) ([]float64, error) {
	return Current().Driver.FloatsP(d.Values[i].([]byte))
}
