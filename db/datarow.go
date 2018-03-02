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

// DataRow row struct
type DataRow map[string]interface{}

// Insert datarow item
func (dr DataRow) Insert(table string) error {
	_, err := NewInsertBuilder(table).Insert(dr)
	return err
}

// TxInsert datarow item
func (dr DataRow) TxInsert(tx *Tx, table string) error {
	_, err := NewInsertBuilder(table).TxInsert(tx, dr)
	return err
}

// ChangedData item
func (dr DataRow) ChangedData(newRow DataRow) DataRow {
	var dat = DataRow{}
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
func (dr DataRow) Update(table string, where string, params ...interface{}) error {
	_, err := NewUpdateBuilder(table).Where(where, params...).Update(dr)
	return err
}

// Save datarow item
// If exists then update
// If not found then insert
func (dr DataRow) Save(table string, where string, params ...interface{}) error {
	if NewExistsBuilder(table).Where(where, params...).Exists() {
		_, err := NewUpdateBuilder(table).Where(where, params...).Update(dr)
		return err
	}
	_, err := NewInsertBuilder(table).Insert(dr)
	return err
}

// InsertIfNotFound inert data
func (dr DataRow) InsertIfNotFound(table string, where string, params ...interface{}) (bool, error) {
	if NewExistsBuilder(table).Where(where, params...).Exists() {
		return true, nil
	}
	_, err := NewInsertBuilder(table).Insert(dr)
	return false, err
}

// Add datarow item
// if current datarow is nil, then create
func (dr DataRow) Add(field string, v interface{}) {
	dr[field] = v
}

// IsEmpty check empty
func (dr DataRow) IsEmpty() bool {
	return len(dr) == 0
}

// Clone func
func (dr DataRow) Clone() DataRow {
	row := DataRow{}
	for k, v := range dr {
		row[k] = v
	}
	return row
}

// IsSet func
func (dr DataRow) IsSet(field string) bool {
	_, ok := dr[field]
	return ok
}

// IsNull func
func (dr DataRow) IsNull(field string) bool {
	if !dr.IsSet(field) {
		return true
	}

	return dr[field] == nil
}

// Fix2JsonData remove NaN, +Inf, -Inf
func (dr DataRow) Fix2JsonData() DataRow {
	for k := range dr {
		v := dr.Float(k)
		if math.IsNaN(v) || math.IsInf(v, 1) || math.IsInf(v, -1) {
			dr[k] = nil
		}
	}
	return dr
}

// Bytes2String convert all bytes to string type
func (dr DataRow) Bytes2String() DataRow {
	for k, v := range dr {
		switch v.(type) {
		case []byte:
			dr[k] = string(v.([]byte))
		}
	}
	return dr
}

// Bytes2NumericByFields convert to number type
func (dr DataRow) Bytes2NumericByFields(fields []string) DataRow {
	for k := range dr {
		if InStrings(fields, k) {
			dr[k] = dr.Float(k)
		}
	}
	return dr
}

func (dr DataRow) Bytes2NumericBySubfix(subfix string) DataRow {
	nn := len(subfix)
	for k := range dr {
		if k[len(k)-nn:] != subfix {
			continue
		}

		dr[k] = dr.Float(k)
	}
	return dr
}

// Bool return
func (dr DataRow) Bool(field string) bool {
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
func (dr DataRow) BoolDefault(field string, v bool) bool {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Bool(field)
}

// Bytes return
func (dr DataRow) Bytes(field string) []byte {
	if dr.IsNull(field) {
		return B_EmptyString
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
		return B_EmptyString
	}
}

// String return
func (dr DataRow) String(field string) string {
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
func (dr DataRow) StringDefault(field string, v string) string {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.String(field)
}

// Uint64 return
func (dr DataRow) Uint64(field string) uint64 {
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
func (dr DataRow) Int64(field string) int64 {
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
			panic(err)
		}
		return i

	default:
		panic(fmt.Sprintf("unknow field %s, can not convert to int64. this field type is %s", field, reflect.TypeOf(dr[field])))
	}

}

// Int64Default bool
func (dr DataRow) Int64Default(field string, v int64) int64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int64(field)
}

// Int return
func (dr DataRow) Int(field string) int {
	return int(dr.Int64(field))
}

// IntDefault bool
func (dr DataRow) IntDefault(field string, v int) int {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int(field)
}

// Float float64
func (dr DataRow) Float(field string) float64 {
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
		panic("unkonw float type to convert")
	}
}

// FloatDefault bool
func (dr DataRow) FloatDefault(field string, v float64) float64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Float(field)
}

// Int64s return []int64
func (dr DataRow) Int64s(field string) []int64 {
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
func (dr DataRow) Int64sDefault(field string, v []int64) []int64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int64s(field)
}

// Ints return []int
func (dr DataRow) Ints(field string) []int {
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
func (dr DataRow) IntsDefault(field string, v []int) []int {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Ints(field)
}

func (dr DataRow) ParseNumberSlice(field string, ptr interface{}) error {
	return Current().Driver.ParseNumberSlice(dr.Bytes(field), ptr)
}

func (dr DataRow) ParseStringSlice(field string, ptr interface{}) error {
	return Current().Driver.ParseStringSlice(dr.Bytes(field), ptr)
}

func (dr DataRow) Floats(field string) []float64 {
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
func (dr DataRow) Strings(field string) []string {
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
		l := len(s)
		if l > 0 && Current().Driver.DriverName() == drivers.DriverPSQL {
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
func (dr DataRow) StringsDefault(field string, v []string) []string {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Strings(field)
}

// DataRow func
func (dr DataRow) DataRow(field string) DataRow {
	switch dr[field].(type) {
	case map[string]interface{}:
		return DataRow(dr[field].(map[string]interface{}))

	case []byte, string:
		v := make(map[string]interface{}, 0)
		dr.JsonParse(field, &v)
		return DataRow(v)

	default:
		panic("DataRow unknow data type")
	}

}

func (dr DataRow) DataSet(field string) DataSet {
	switch dr[field].(type) {
	case []DataRow:
		return DataSet(dr[field].([]DataRow))

	case []byte, string:
		v := make([]DataRow, 0)
		dr.JsonParse(field, &v)
		return DataSet(v)

	default:
		panic("DataSet unknow data type")
	}

}

func (dr DataRow) Hstore(field string) map[string]string {
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

func (dr DataRow) JsonParse(field string, v interface{}) error {
	if err := json.Unmarshal(dr.Bytes(field), v); err != nil {
		return err
	}
	return nil
}

func (dr DataRow) Time(field string) time.Time {
	if dr.IsNull(field) {
		return time.Unix(0, 0)
	}
	switch dr[field].(type) {
	case string:
		format := DBTimeFormat
		if len(dr[field].(string)) == 10 {
			format = "2006-01-02"
		}
		t, err := time.Parse(format, dr[field].(string))
		if err != nil {
			panic(err)
		}
		return t

	case time.Time:
		return dr[field].(time.Time)

	case []byte:
		str := string(dr[field].([]byte))
		format := DBTimeFormat
		if len(str) == 10 {
			format = "2006-01-02"
		}
		t, err := time.Parse(format, str)
		if err != nil {
			panic(err)
		}
		return t

	default:
		panic(fmt.Sprintf("StringSlice unknow data type %s, field %s", reflect.TypeOf(dr[field]).String(), field))
	}
}

// CopyTo func
func (dr DataRow) CopyTo(vo interface{}) error {
	src, err := json.Marshal(dr)
	if err != nil {
		return err
	}

	return json.Unmarshal(src, vo)
}

// ConvertTo func
func (dr DataRow) ConvertTo(vo interface{}) error {
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

	var field string
	var sf reflect.StructField

	n := typ.NumField()
	for i := 0; i < n; i++ {
		sf = typ.Field(i)

		field = sf.Tag.Get("json")

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
				err := dr.JsonParse(field, v.Interface())
				if err != nil {
					return err
				}
				val.Field(i).Set(v)

			default:
				val.Field(i).Set(reflect.ValueOf(dr[field]))

			}

		case reflect.Struct, reflect.Interface:
			switch sf.Type.String() {
			case "time.Time":
				val.Field(i).Set(reflect.ValueOf(dr.Time(field)))

			default:
				switch dr[field].(type) {
				case []byte, string:
					// data type is json
					v := reflect.New(sf.Type)
					err := dr.JsonParse(field, v.Interface())
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

			case "map[string]interface {}", "db.DataRow":
				val.Field(i).Set(reflect.ValueOf(dr.DataRow(field)))

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

		defer func() {
			if err := recover(); err != nil {
				panic(fmt.Sprint("Field(", i, ") ", field))
			}
		}()
	}
	return nil
}

// Split2Slice split map data
func (dr DataRow) Split2Slice() ([]string, []interface{}) {
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
func (dr DataRow) Merge(row DataRow) DataRow {
	for k, v := range row {
		dr[k] = v
	}

	return dr
}

// Keys list
func (dr DataRow) Keys() []string {
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
func (dr DataRow) Len() int {
	return len(dr)
}
