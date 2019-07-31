package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kere/gno/db/drivers"
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

// Add datarow item
// if current datarow is nil, then create
func (dr MapRow) Add(field string, v interface{}) {
	dr[field] = v
}

// IsEmpty check empty
func (dr MapRow) IsEmpty() bool {
	return len(dr) == 0
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
			dr[k] = string(v.([]byte))
		}
	}
	return dr
}

// // Bytes2NumericByFields convert to number type
// func (dr MapRow) Bytes2NumericByFields(fields []string) MapRow {
// 	for k := range dr {
// 		if InStrings(fields, k) {
// 			dr[k] = dr.Float(k)
// 		}
// 	}
// 	return dr
// }
//
// // Bytes2NumericBySubfix func
// func (dr MapRow) Bytes2NumericBySubfix(subfix string) MapRow {
// 	nn := len(subfix)
// 	for k := range dr {
// 		if k[len(k)-nn:] != subfix {
// 			continue
// 		}
//
// 		dr[k] = dr.Float(k)
// 	}
// 	return dr
// }

// Bool return
func (dr MapRow) Bool(field string) bool {
	switch dr[field].(type) {
	case int, int64, int32, int16, int8, float32, float64, uint, uint64, uint32, uint16, uint8:
		return dr.Int64(field) > 0

	case string:
		switch dr[field].(string) {
		case "t", "T":
			return true

		default:
			return false

		}
	case bool:
		return dr[field].(bool)

	default:
		return false
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

	case int, int64, int32, int8:
		return []byte(strconv.FormatInt(dr.Int64(field), 10))

	case uint, uint64, uint32, uint8:
		return []byte(strconv.FormatUint(dr.Uint64(field), 10))

	default:
		return BEmptyString
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
		return string(dr[field].([]byte))

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

// Uint64 return
func (dr MapRow) Uint64(field string) uint64 {
	if dr[field] == nil {
		return 0
	}
	// return dr[field].(int64)
	switch dr[field].(type) {
	case int, int64, int32, int16, int8:
		return uint64(dr.Int64(field))

	case float64, float32:
		return uint64(dr.Float(field))

	case uint:
		return uint64(dr[field].(uint))

	case uint64:
		return dr[field].(uint64)

	case uint32:
		return uint64(dr[field].(uint32))

	case uint16:
		return uint64(dr[field].(uint16))

	case uint8:
		return uint64(dr[field].(uint8))

	case string:
		i, err := strconv.ParseInt(dr[field].(string), 10, 64)
		if err != nil {
			panic(err)
		}
		return uint64(i)

	default:
		panic(fmt.Sprintf("unknow field %s, can not convert to int64. this field type is %s", field, reflect.TypeOf(dr[field])))
	}

}

// Int64 return
func (dr MapRow) Int64(field string) int64 {
	if dr[field] == nil {
		return 0
	}
	// return dr[field].(int64)
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

	case int16:
		return int64(dr[field].(int16))

	case int8:
		return int64(dr[field].(int8))

	case uint, uint64, uint32, uint16, uint8:
		return int64(dr.Uint64(field))

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
	case uint, uint64, uint32, uint16, uint8:
		return float64(dr.Uint64(field))

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

	case string, []byte:
		v := make([]int64, 0)
		if err := dr.ParseNumberSlice(field, &v); err != nil {
			panic(err)
		}
		return v

	default:
		panic("Int64Slice unknow data type")
	}

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

	case string, []byte:
		v := make([]int, 0)
		if err := dr.ParseNumberSlice(field, &v); err != nil {
			panic(err)
		}
		return v

	default:
		panic("Int64Slice unknow data type")
	}

}

// IntsDefault bool
func (dr MapRow) IntsDefault(field string, v []int) []int {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Ints(field)
}

//ParseNumberSlice err
func (dr MapRow) ParseNumberSlice(field string, ptr interface{}) error {
	return Current().Driver.ParseNumberSlice(dr.Bytes(field), ptr)
}

//ParseStringSlice err
func (dr MapRow) ParseStringSlice(field string, ptr interface{}) error {
	return Current().Driver.ParseStringSlice(dr.Bytes(field), ptr)
}

// Floats []float64
func (dr MapRow) Floats(field string) []float64 {
	switch dr[field].(type) {
	case []float64:
		return dr[field].([]float64)

	case string, []byte:
		b := dr.Bytes(field)
		if len(b) > 0 {
			s := make([]float64, 0)
			if err := Current().Driver.ParseNumberSlice(b, &s); err != nil {
				panic(err)
			}
			return s
		}

		return []float64{}

	default:
		panic("Floats unknow data type")
	}

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
	case string, []byte:

		s := dr.String(field)
		if s == "" {
			return []string{}
		}

		l := len(s)
		if l > 0 && Current().Driver.Name() == drivers.DriverPSQL {
			if s[:1] == "{" && s[l-1:l] == "}" {
				if arr, err := Current().Driver.StringSlice([]byte(s)); err == nil {
					return arr
				}
				return []string{}
			}
		}
		return strings.Split(s, ",")
	default:
		return []string{}
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
		dr.JSONParse(field, &v)
		return MapRow(v)

	default:
		panic("MapRow unknow data type")
	}

}

//MapRows func
func (dr MapRow) MapRows(field string) MapRows {
	switch dr[field].(type) {
	case []MapRow:
		return MapRows(dr[field].([]MapRow))

	case []byte, string:
		v := make([]MapRow, 0)
		dr.JSONParse(field, &v)
		return MapRows(v)

	default:
		panic("MapRows unknow data type")
	}

}

// Hstore dr
func (dr MapRow) Hstore(field string) map[string]string {
	switch dr[field].(type) {
	case map[string]string:
		return dr[field].(map[string]string)

	case []byte, string:
		data, err := Current().Driver.HStore(dr.Bytes(field))
		if err != nil {
			panic(err)
		}
		return data

	default:
		panic("Hstore unknow data type")
	}

}

// JSONParse data row to parse jSON field
func (dr MapRow) JSONParse(field string, v interface{}) error {
	return json.Unmarshal(dr.Bytes(field), v)
}

// Time default
func (dr MapRow) Time(field string) time.Time {
	return dr.TimeParse(field, DateTimeFormat)
}

// TimeParse default
func (dr MapRow) TimeParse(field, layout string) time.Time {
	if dr.IsNull(field) {
		return time.Unix(0, 0)
	}

	switch dr[field].(type) {
	case string:
		t, err := time.Parse(layout, dr[field].(string))
		// t, err := time.ParseInLocation(format, dr[field].(string), loc)
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

// CopyToVO func
func (dr MapRow) CopyToVO(vo interface{}) error {
	typ := reflect.TypeOf(vo)
	if typ.Kind() != reflect.Ptr {
		return errors.New("arg vo must be a kind of ptr")
	}

	typ = typ.Elem()
	val := reflect.ValueOf(vo)
	if val.IsNil() {
		return errors.New("copy to struct failed, vo is nil")
	}

	val = val.Elem()
	if !val.IsValid() {
		return errors.New("copy to struct failed, vo is invalid")
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("MapRow ConvertTo failed ")
		}
	}()

	var field string
	var sf reflect.StructField

	n := typ.NumField()
	for i := 0; i < n; i++ {
		sf = typ.Field(i)

		field = sf.Tag.Get(FieldJSON)

		if field == "" {
			continue
		}

		if !dr.IsSet(field) {
			continue
		}

		// fmt.Println(sf.Name, sf.Type.Kind())
		switch sf.Type.Kind() {
		case reflect.Ptr:
			switch dr[field].(type) {
			case []byte, string:
				// data type is json
				v := reflect.New(sf.Type.Elem())
				err := dr.JSONParse(field, v.Interface())
				if err != nil {
					return err
				}
				val.Field(i).Set(v)

			default:
				val.Field(i).Set(reflect.ValueOf(dr[field]))

			}

		case reflect.Struct, reflect.Interface:
			switch sf.Type.String() {
			case timeClassName:
				val.Field(i).Set(reflect.ValueOf(dr.Time(field)))

			default:
				switch dr[field].(type) {
				case []byte, string:
					// data type is json
					v := reflect.New(sf.Type)
					err := dr.JSONParse(field, v.Interface())
					if err != nil {
						return err
					}
					val.Field(i).Set(v.Elem())

				default:
					val.Field(i).Set(reflect.ValueOf(dr[field]))

				}
			}

		case reflect.Map:
			switch sf.Type.String() {
			case "map[string]string":
				val.Field(i).Set(reflect.ValueOf(dr.Hstore(field)))

			case "db.MapRow", "map[string]interface {}":
				val.Field(i).Set(reflect.ValueOf(dr.MapRow(field)))

			default:
				return fmt.Errorf("unkonw map data type of %s", sf.Type.String())
			}

		case reflect.Slice, reflect.Array:
			switch arrayBaseType(sf.Type).Kind() {
			case reflect.String:
				// array base type is string
				switch dr[field].(type) {
				case []byte, string:
					cv := reflect.New(val.Field(i).Type())
					err := Current().Driver.ParseStringSlice(dr.Bytes(field), cv.Interface())
					if err != nil {
						return err
					}
					val.Field(i).Set(cv.Elem())

				default:
					b, err := json.Marshal(dr[field])
					if err != nil {
						return err
					}
					cv := reflect.New(val.Field(i).Type())
					err = Current().Driver.ParseStringSlice(b, cv.Interface())
					if err != nil {
						return err
					}
					val.Field(i).Set(cv.Elem())
				}

			case reflect.Int64, reflect.Int, reflect.Int32, reflect.Float64, reflect.Float32:
				// array base type is number
				switch dr[field].(type) {
				case []byte, string:
					cv := reflect.New(val.Field(i).Type())
					err := Current().Driver.ParseNumberSlice(dr.Bytes(field), cv.Interface())
					if err != nil {
						return err
					}
					val.Field(i).Set(cv.Elem())

				default:
					b, err := json.Marshal(dr[field])
					if err != nil {
						return err
					}
					cv := reflect.New(val.Field(i).Type())
					err = Current().Driver.ParseNumberSlice(b, cv.Interface())
					if err != nil {
						return err
					}
					val.Field(i).Set(cv.Elem())
				}

			}

		case reflect.String:
			val.Field(i).SetString(dr.String(field))

		case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			val.Field(i).SetInt(dr.Int64(field))

		case reflect.Uint64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			val.Field(i).SetUint(uint64(dr.Int64(field)))

		case reflect.Bool:
			val.Field(i).SetBool(dr.Bool(field))

		case reflect.Float64, reflect.Float32:
			val.Field(i).SetFloat(dr.Float(field))
		}
	}
	return nil
}

// SplitData split map data
func (dr MapRow) SplitData() ([]string, []interface{}) {
	l := len(dr)
	keys := make([]string, l)
	vals := make([]interface{}, l)
	i := 0
	for k := range dr {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	for i := 0; i < l; i++ {
		vals[i] = dr[keys[i]]
	}

	return keys, vals
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
