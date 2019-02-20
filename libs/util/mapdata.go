package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// MapData util mapdata
type MapData map[string]interface{}

// Empty check empty
func (dr MapData) Empty() bool {
	return len(dr) == 0
}

// IsSet isset
func (dr MapData) IsSet(field string) bool {
	_, ok := dr[field]
	return ok
}

// IsNull nil
func (dr MapData) IsNull(field string) bool {
	return dr[field] == nil
}

//Bytes2String convert all to string if type is []byte
func (dr MapData) Bytes2String() MapData {
	for k, v := range dr {

		switch v.(type) {
		case []byte:
			dr[k] = string(v.([]byte))
		}
	}

	return dr
}

// Bool bool
func (dr MapData) Bool(field string) bool {
	switch dr[field].(type) {
	case int64, int32, int16, int8, int, float32, float64:
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
func (dr MapData) BoolDefault(field string, v bool) bool {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Bool(field)
}

// Bytes return
func (dr MapData) Bytes(field string) []byte {
	if dr.IsNull(field) {
		return []byte("")
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
		return []byte("")
	}
}

// BytesDefault bool
func (dr MapData) BytesDefault(field string, v []byte) []byte {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Bytes(field)
}

// Rune return
func (dr MapData) Rune(field string) rune {
	str := dr.String(field)
	if str == "" {
		return 0
	}
	return rune(str[0])
}

// String string
func (dr MapData) String(field string) string {
	switch dr[field].(type) {
	case string:
		return dr[field].(string)

	case []byte:
		return string(dr[field].([]byte))

	case nil:
		return ""

	default:
		return fmt.Sprint(dr[field])
	}
}

// StringDefault bool
func (dr MapData) StringDefault(field string, v string) string {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.String(field)
}

// Uint64 return
func (dr MapData) Uint64(field string) uint64 {
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
func (dr MapData) Int64(field string) int64 {
	// return dr[field].(int64)
	switch dr[field].(type) {
	case float64, float32:
		return int64(dr.Float(field))

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
func (dr MapData) Int64Default(field string, v int64) int64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int64(field)
}

// Int int
func (dr MapData) Int(field string) int {
	return int(dr.Int64(field))
}

// IntDefault bool
func (dr MapData) IntDefault(field string, v int) int {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int(field)
}

// Int64s []int64
func (dr MapData) Int64s(field string) []int64 {
	switch dr[field].(type) {
	case []int64:
		return dr[field].([]int64)

	case []interface{}:
		vals := dr[field].([]interface{})
		v := make([]int64, len(vals))
		for i, val := range vals {
			switch val.(type) {
			case int64:
				v[i] = val.(int64)
			case int:
				v[i] = int64(val.(int))
			case float64:
				v[i] = int64(val.(float64))
			case string:
				n, err := strconv.ParseInt(val.(string), 10, 64)
				if err == nil {
					v[i] = n
				}
			default:
				panic("unknow type in function Ints:" + reflect.TypeOf(val).String())
			}
		}

		return v

	default:
		return []int64{}
	}
}

// Int64sDefault bool
func (dr MapData) Int64sDefault(field string, v []int64) []int64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Int64s(field)
}

// Ints []int
func (dr MapData) Ints(field string) []int {
	switch dr[field].(type) {
	case []int:
		return dr[field].([]int)

	case []interface{}:
		vals := dr[field].([]interface{})
		v := make([]int, len(vals))
		for i, val := range vals {
			switch val.(type) {
			case int64:
				v[i] = int(val.(int64))
			case int:
				v[i] = val.(int)
			case float64:
				v[i] = int(val.(float64))
			case string:
				n, err := strconv.ParseInt(val.(string), 10, 64)
				if err == nil {
					v[i] = int(n)
				}
			default:
				panic("unknow type in function Ints:" + reflect.TypeOf(val).String())
			}
		}

		return v
	case []byte:
		s := dr[field].([]byte)
		var v []int
		if s[0] == '[' {
			json.Unmarshal(s, &v)
		}

		return v

	case string:
		s := dr[field].(string)
		var v []int
		if s[0] == '[' {
			json.Unmarshal([]byte(s), &v)
		}

		return v

	default:
		return []int{}
	}
}

// IntsDefault bool
func (dr MapData) IntsDefault(field string, v []int) []int {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Ints(field)
}

// Strings []string
func (dr MapData) Strings(field string) []string {
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
		return strings.Split(s, ",")
	default:
		return []string{}
	}
}

// StringsDefault bool
func (dr MapData) StringsDefault(field string, v []string) []string {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Strings(field)
}

// Float float64
func (dr MapData) Float(field string) float64 {
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
func (dr MapData) FloatDefault(field string, v float64) float64 {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Float(field)
}

// MapDatas []MapData
func (dr MapData) MapDatas(field string) []MapData {
	switch dr[field].(type) {
	case []interface{}:
		data := dr[field].([]interface{})
		items := make([]MapData, len(data))
		for i := range data {
			items[i] = (MapData)(data[i].(map[string]interface{}))
		}
		return items
	case []MapData:
		return dr[field].([]MapData)
	default:
		return nil
	}
}

// MapData mapData
func (dr MapData) MapData(field string) MapData {
	switch dr[field].(type) {
	case map[string]interface{}:
		return MapData(dr[field].(map[string]interface{}))

	case []byte, string:
		v := make(map[string]interface{}, 0)
		dr.JSONParse(field, &v)
		return MapData(v)

	default:
		panic("Hstore unknow data type")
	}
}

// MapDataDefault bool
func (dr MapData) MapDataDefault(field string, v MapData) MapData {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.MapData(field)
}

// ParseTo json object
func (dr MapData) ParseTo(field string, to interface{}) error {
	b, err := json.Marshal(dr[field])
	if err != nil {
		return err
	}

	return json.Unmarshal(b, to)
}

// JSONParse parse json
func (dr MapData) JSONParse(field string, v interface{}) error {
	if err := json.Unmarshal(dr.Bytes(field), v); err != nil {
		return err
	}
	return nil
}

// CopyTo vo
func (dr MapData) CopyTo(v interface{}) error {
	src, err := json.Marshal(dr)
	if err != nil {
		return err
	}
	return json.Unmarshal(src, v)
}

// Time time
func (dr MapData) Time(field string) time.Time {
	if dr.IsNull(field) {
		return time.Unix(0, 0)
	}
	switch dr[field].(type) {
	case string:
		format := DBTimeFormat
		if len(dr[field].(string)) == 10 {
			format = "2006-01-02"
		}
		t, err := time.ParseInLocation(format, dr[field].(string), time.Local)
		if err != nil {
			panic(err)
		}
		return t

	case time.Time:
		return dr[field].(time.Time)

	default:
		panic("Time unknow data type")
	}
}

// TimeDefault bool
func (dr MapData) TimeDefault(field string, v time.Time) time.Time {
	if !dr.IsSet(field) || dr.IsNull(field) {
		return v
	}
	return dr.Time(field)
}

// Clone data
func (dr MapData) Clone() MapData {
	row := MapData{}
	for k, v := range dr {
		row[k] = v
	}
	return row
}

// ArgMult func
func (dr MapData) ArgMult(field string, v float64) MapData {
	if !dr.IsSet(field) {
		return dr
	}

	dr[field] = dr.Float(field) * v
	return dr
}

// ArgPlus func
func (dr MapData) ArgPlus(field string, v float64) MapData {
	if !dr.IsSet(field) {
		return dr
	}

	dr[field] = dr.Float(field) + v
	return dr
}

// ArgDiv func
func (dr MapData) ArgDiv(field string, v float64) MapData {
	if !dr.IsSet(field) {
		return dr
	}

	dr[field] = dr.Float(field) / v
	return dr
}
