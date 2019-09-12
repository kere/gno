package db

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/kere/gno/libs/util"
)

// MapRow row struct
type MapRow map[string]interface{}

// Insert datarow item
func (dr MapRow) Insert(table string) error {
	ins := InsertBuilder{}
	_, err := ins.Table(table).Insert(dr)
	return err
}

// TxInsert datarow item
func (dr MapRow) TxInsert(tx *Tx, table string) error {
	ins := InsertBuilder{}
	_, err := ins.Table(table).TxInsert(tx, dr)
	return err
}

// ChangedData item
func (dr MapRow) ChangedData(newRow MapRow) MapRow {
	var dat = MapRow{}
	var val interface{}
	var isok bool
	for k, v := range dr {
		if val, isok = newRow[k]; !isok {
			continue
		}

		typ := reflect.TypeOf(v)
		switch typ.Kind() {
		default:
			if v == val {
				continue
			}
		case reflect.Struct:
		case reflect.Map:

		case reflect.Uint8:
			if dr.String(k) == newRow.String(k) {
				continue
			}
		case reflect.Int, reflect.Int64, reflect.Int16, reflect.Int8:
			if dr.Int64(k) == newRow.Int64(k) {
				continue
			}
		case reflect.Float32, reflect.Float64:
			if dr.Float(k) == newRow.Float(k) {
				continue
			}
		case reflect.Slice:
			vv := reflect.ValueOf(v)
			vvNew := reflect.ValueOf(newRow[k])
			if vv.Kind() == reflect.Ptr {
				vv = vv.Elem()
			}
			l := vv.Len()
			if l > 0 && l == vvNew.Len() {
				isEq := true
				for i := 0; i < l; i++ {
					if vv.Index(i).Interface() != vvNew.Index(i).Interface() {
						isEq = false
						break
					}
				}
				if isEq {
					continue
				}
			}

		}

		dat[k] = newRow[k]
	}

	return dat
}

// Update datarow item
func (dr MapRow) Update(table string, where string, params ...interface{}) error {
	u := UpdateBuilder{table: table}
	_, err := u.Where(where, params...).Update(dr)
	return err
}

// Save datarow item
// If exists then update
// If not found then insert
func (dr MapRow) Save(table string, where string, params ...interface{}) error {
	e := ExistsBuilder{}
	if e.Table(table).Where(where, params...).Exists() {
		e := UpdateBuilder{table: table}
		_, err := e.Where(where, params...).Update(dr)
		return err
	}
	ins := InsertBuilder{}
	_, err := ins.Table(table).Insert(dr)
	return err
}

// InsertIfNotFound inert data
func (dr MapRow) InsertIfNotFound(table string, where string, params ...interface{}) (bool, error) {
	e := ExistsBuilder{}
	if e.Table(table).Where(where, params...).Exists() {
		return true, nil
	}
	ins := InsertBuilder{}
	_, err := ins.Table(table).Insert(dr)
	return false, err
}

// IsEmpty check empty
func (dr MapRow) IsEmpty() bool {
	return len(dr) == 0
}

// IsEmptyValue check empty
func (dr MapRow) IsEmptyValue(field string) bool {
	if dr.IsNull(field) {
		return true
	}
	return isEmptyValue(reflect.ValueOf(dr[field]))
}

// Clone func
func (dr MapRow) Clone() MapRow {
	row := MapRow{}
	for k, v := range dr {
		row[k] = v
	}
	return row
}

// IsSet func
func (dr MapRow) IsSet(field string) bool {
	_, ok := dr[field]
	return ok
}

// IsNull func
func (dr MapRow) IsNull(field string) bool {
	if !dr.IsSet(field) {
		return true
	}

	return dr[field] == nil
}

// Fix2JsonData remove NaN, +Inf, -Inf
func (dr MapRow) Fix2JsonData() MapRow {
	for k := range dr {
		v := dr.Float(k)
		if math.IsNaN(v) || math.IsInf(v, 1) || math.IsInf(v, -1) {
			dr[k] = nil
		}
	}
	return dr
}

// Bytes2String convert all bytes to string type
func (dr MapRow) Bytes2String() MapRow {
	for k, v := range dr {
		switch v.(type) {
		case []byte:
			if v == nil {
				continue
			}
			dr[k] = util.Bytes2Str(v.([]byte))
		}
	}
	return dr
}

// Bool return
func (dr MapRow) Bool(field string) bool {
	switch dr[field].(type) {
	case int, int64:
		return dr.Int64(field) > 0

	case string:
		v := dr[field].(string)
		if len(v) != 1 {
			return false
		}
		return v[0] == 't' || v[0] == 'T'
	case bool:
		return dr[field].(bool)

	default:
		panic(ErrType)
	}
}

// BoolDefault bool
func (dr MapRow) BoolDefault(field string, v bool) bool {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Bool(field)
}

// Bytes return
func (dr MapRow) Bytes(field string) []byte {
	if dr.IsNull(field) {
		return BEmptyString
	}

	switch dr[field].(type) {
	case string:
		return []byte(dr[field].(string))

	case []byte:
		return dr[field].([]byte)

	case float64, float32:
		return []byte(strconv.FormatFloat(dr.Float(field), 'f', -1, 64))

	case int, int64:
		return []byte(strconv.FormatInt(dr.Int64(field), 10))

	default:
		src, err := json.Marshal(dr[field])
		if err != nil {
			panic(err)
		}
		return src
	}
}

// Rune return
func (dr MapRow) Rune(field string) rune {
	str := dr.String(field)
	if str == "" {
		return 0
	}
	return rune(str[0])
}

// String return
func (dr MapRow) String(field string) string {
	if dr.IsNull(field) {
		return ""
	}

	switch dr[field].(type) {
	case string:
		return dr[field].(string)

	case []byte:
		return util.Bytes2Str(dr[field].([]byte))

	default:
		return fmt.Sprint(dr[field])
	}

}

// StringDefault bool
func (dr MapRow) StringDefault(field string, v string) string {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.String(field)
}

// Int64 return
func (dr MapRow) Int64(field string) int64 {
	if dr[field] == nil {
		return 0
	}

	switch dr[field].(type) {
	case float64, float32:
		return int64(dr.Float(field))

	case bool:
		if dr[field].(bool) {
			return 1
		}
		return 0

	case int:
		return int64(dr[field].(int))

	case int64:
		return dr[field].(int64)

	case int32:
		return int64(dr[field].(int32))

	case string:
		i, err := strconv.ParseInt(dr[field].(string), 10, 64)
		if err != nil {
			v, err := strconv.ParseFloat(dr[field].(string), 64)
			if err != nil {
				panic(err)
			}
			i = int64(v)
		}
		return i

	default:
		panic(fmt.Sprintf("unknow field %s, can not convert to int64. this field type is %s", field, reflect.TypeOf(dr[field])))
	}

}

// Int64Default bool
func (dr MapRow) Int64Default(field string, v int64) int64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int64(field)
}

// Int return
func (dr MapRow) Int(field string) int {
	if dr[field] == nil {
		return 0
	}
	return int(dr.Int64(field))
}

// IntDefault bool
func (dr MapRow) IntDefault(field string, v int) int {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int(field)
}

// Float float64
func (dr MapRow) Float(field string) float64 {
	if dr[field] == nil {
		return 0
	}

	switch dr[field].(type) {
	case []byte:
		f, err := strconv.ParseFloat(string(dr[field].([]byte)), 64)
		if err != nil {
			panic(err)
		}
		return f
	case string:
		f, err := strconv.ParseFloat(dr[field].(string), 64)
		if err != nil {
			panic(err)
		}
		return f

	case float64:
		return dr[field].(float64)

	case float32:
		return float64(dr[field].(float32))

	case int, int64, int32, int16, int8:
		return float64(dr.Int64(field))

	default:
		panic("unkonw float type to convert, maybe is nil[" + field + "]")
	}
}

// FloatDefault bool
func (dr MapRow) FloatDefault(field string, v float64) float64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Float(field)
}

// Int64s return []int64
func (dr MapRow) Int64s(field string) []int64 {
	switch dr[field].(type) {
	case []int64:
		return dr[field].([]int64)
	case []byte:
		v, err := Current().Driver.Int64s(dr[field].([]byte))
		if err != nil {
			panic(err)
		}
		return v

	case string:
		s := dr[field].(string)
		var v []int64
		if err := Current().Driver.ParseNumberSlice(util.Str2Bytes(s), &v); err != nil {
			panic(err)
		}
		return v

	case []interface{}:
		vals := dr[field].([]interface{})
		l := len(vals)
		arr := make([]int64, l)
		if l == 0 {
			return arr
		}

		typ := reflect.TypeOf(vals[0])
		switch typ.Kind() {
		case reflect.Int:
			for i := 0; i < l; i++ {
				arr[i] = int64(vals[i].(int))
			}
		case reflect.Int32:
			for i := 0; i < l; i++ {
				arr[i] = int64(vals[i].(int32))
			}
		case reflect.Int64:
			for i := 0; i < l; i++ {
				arr[i] = vals[i].(int64)
			}
		case reflect.Float64:
			for i := 0; i < l; i++ {
				arr[i] = int64(vals[i].(float64))
			}
		case reflect.Float32:
			for i := 0; i < l; i++ {
				arr[i] = int64(vals[i].(float32))
			}
		default:
			panic("Int64Slice unknow data type:[]" + typ.String())
		}
		return arr
	}

	panic("Int64Slice unknow data type:" + reflect.TypeOf(dr[field]).String())
}

// Int64sDefault bool
func (dr MapRow) Int64sDefault(field string, v []int64) []int64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int64s(field)
}

// Ints return []int
func (dr MapRow) Ints(field string) []int {
	if dr.IsNull(field) {
		return []int{}
	}

	switch dr[field].(type) {
	case []int:
		return dr[field].([]int)

	case []byte:
		v := dr[field].([]byte)
		var val []int
		if err := Current().Driver.ParseNumberSlice(v, &val); err != nil {
			panic(err)
		}
		return val
	case string:
		v := dr[field].(string)
		var val []int
		if err := Current().Driver.ParseNumberSlice(util.Str2Bytes(v), &val); err != nil {
			panic(err)
		}
		return val

	case []interface{}:
		vals := dr[field].([]interface{})
		l := len(vals)
		arr := make([]int, l)
		if l == 0 {
			return arr
		}

		typ := reflect.TypeOf(vals[0])
		switch typ.Kind() {
		case reflect.Int:
			for i := 0; i < l; i++ {
				arr[i] = vals[i].(int)
			}
		case reflect.Int32:
			for i := 0; i < l; i++ {
				arr[i] = int(vals[i].(int32))
			}
		case reflect.Int64:
			for i := 0; i < l; i++ {
				arr[i] = int(vals[i].(int64))
			}
		case reflect.Float64:
			for i := 0; i < l; i++ {
				arr[i] = int(vals[i].(float64))
			}
		case reflect.Float32:
			for i := 0; i < l; i++ {
				arr[i] = int(vals[i].(float32))
			}
		default:
			panic("Int64Slice unknow data type:[]" + typ.String())
		}
		return arr

	}
	panic("Int64Slice unknow data type:" + reflect.TypeOf(dr[field]).String())
}

// IntsDefault bool
func (dr MapRow) IntsDefault(field string, v []int) []int {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Ints(field)
}

// Floats []float64
func (dr MapRow) Floats(field string) []float64 {
	switch dr[field].(type) {
	case []float64:
		return dr[field].([]float64)

	case []byte:
		v := dr[field].([]byte)
		val, err := Current().Driver.Floats(v)
		if err != nil {
			panic(err)
		}
		return val

	case string:
		b := dr.Bytes(field)
		if len(b) > 0 {
			var s []float64
			if err := Current().Driver.ParseNumberSlice(b, &s); err != nil {
				panic(err)
			}
			return s
		}
	case []interface{}:
		vals := dr[field].([]interface{})
		l := len(vals)
		arr := make([]float64, l)
		if l == 0 {
			return arr
		}

		typ := reflect.TypeOf(vals[0])
		switch typ.Kind() {
		case reflect.Int:
			for i := 0; i < l; i++ {
				arr[i] = float64(vals[i].(int))
			}
		case reflect.Int32:
			for i := 0; i < l; i++ {
				arr[i] = float64(vals[i].(int32))
			}
		case reflect.Int64:
			for i := 0; i < l; i++ {
				arr[i] = float64(vals[i].(int64))
			}
		case reflect.Float64:
			for i := 0; i < l; i++ {
				arr[i] = vals[i].(float64)
			}
		case reflect.Float32:
			for i := 0; i < l; i++ {
				arr[i] = float64(vals[i].(float32))
			}

		default:
			panic("FloatSlice unknow data type:[]" + typ.String())
		}

		return arr
	}
	panic("FloatSlice unknow data type:" + reflect.TypeOf(dr[field]).String())
}

// Strings []string
func (dr MapRow) Strings(field string) []string {
	if dr.IsNull(field) {
		return []string{}
	}

	switch dr[field].(type) {
	case []string:
		return dr[field].([]string)

	case []interface{}:
		vals := dr[field].([]interface{})
		v := make([]string, len(vals))
		for i, val := range vals {
			v[i] = fmt.Sprint(val)
		}
		return v
	case []byte:
		v := dr[field].([]byte)
		val, err := Current().Driver.Strings(v)
		if err != nil {
			panic(err)
		}
		return val

	case string:
		// {1,1} postgres array => []byte
		s := dr[field].(string)
		if s == "" {
			return nil
		}
		var val []string
		Current().Driver.ParseStringSlice(util.Str2Bytes(s), &val)
		return val

	default:
		panic(ErrType)
	}
}

// StringsDefault bool
func (dr MapRow) StringsDefault(field string, v []string) []string {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Strings(field)
}

// MapRow func
func (dr MapRow) MapRow(field string) MapRow {
	switch dr[field].(type) {
	case map[string]interface{}:
		return MapRow(dr[field].(map[string]interface{}))

	case []byte, string:
		v := make(map[string]interface{}, 0)
		err := json.Unmarshal(dr.Bytes(field), &v)
		if err != nil {
			panic("MapRow unknow data type:" + err.Error())
		}
		return MapRow(v)

	}

	panic("MapRow unknow data type:" + reflect.TypeOf(dr[field]).String())
}

//MapRows func
func (dr MapRow) MapRows(field string) MapRows {
	switch dr[field].(type) {
	case []MapRow:
		return MapRows(dr[field].([]MapRow))

	case []byte, string:
		v := make([]MapRow, 0)
		err := json.Unmarshal(dr.Bytes(field), &v)
		if err != nil {
			panic("MapRows unknow data type:" + err.Error())
		}
		return MapRows(v)

	}

	panic("MapRows unknow data type:" + reflect.TypeOf(dr[field]).String())
}

// // JSONParse data row to parse jSON field
// func (dr MapRow) JSONParse(field string, v interface{}) error {
// 	return json.Unmarshal(dr.Bytes(field), v)
// }

// Time default
func (dr MapRow) Time(field string) time.Time {
	if dr.IsNull(field) {
		return EmptyTime
	}

	if dr.IsEmpty() {
		return EmptyTime
	}
	return dr.TimeParse(field, DateTimeFormat)
}

// TimeParse default
func (dr MapRow) TimeParse(field, layout string) time.Time {
	if dr.IsNull(field) {
		return EmptyTime
	}

	switch dr[field].(type) {
	case string:
		t, err := time.Parse(layout, dr[field].(string))
		if err != nil {
			panic(err)
		}
		return t

	case time.Time:
		return dr[field].(time.Time)

	case []byte:
		t, err := time.Parse(layout, string(dr[field].([]byte)))
		// t, err := time.ParseInLocation(format, str, loc)
		if err != nil {
			panic(err)
		}
		return t

	default:
		panic(fmt.Sprintf("StringSlice unknow data type %s, field %s", reflect.TypeOf(dr[field]).String(), field))
	}
}

// CopyToWithJSON copy to json vo
func (dr MapRow) CopyToWithJSON(vo interface{}) error {
	src, err := json.Marshal(dr)
	if err != nil {
		return err
	}

	return json.Unmarshal(src, vo)
}

// Merge data
func (dr MapRow) Merge(row MapRow) MapRow {
	for k, v := range row {
		dr[k] = v
	}

	return dr
}

// Keys list
func (dr MapRow) Keys() []string {
	l := len(dr)
	arr := make([]string, l)
	i := 0
	for k := range dr {
		arr[i] = k
		i++
	}
	return arr
}

// Len data
func (dr MapRow) Len() int {
	return len(dr)
}
