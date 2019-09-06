package db

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/kere/gno/libs/util"
)

const (
	timeType = "time.Time"
	// tagSkip       = "skip"
	tagJSON      = "json"
	tagSkipEmpty = "skipempty"
	tagType      = "type"
	tagNull      = "null" // 如果isEmpty 则写入null
	tagFormat    = "format"

	typAutoTime = "autotime"
	typPlus     = "plus"
	typPlus1    = "plus1"
	// typVersion      = "version"
	vAll            = "all"
	skipFieldBaseVO = "BaseVO"
	skipFieldTable  = "Table"
)

// StructField class
type StructField struct {
	Field string
	Name  string
}

// Convert class
// Example: sc := NewStructConvert(UserVO{})
// dataset, err := db.NewQueryBuilder("users").Where("id=?",1).Query()
// result := cs.DataSet2Struct(maprow)
// userVO := result[0].(*UserVO)
type Convert struct {
	typ      reflect.Type
	val      reflect.Value
	target   interface{}
	fields   []*StructField
	dbfields [][]byte
}

// NewConvert func
func NewConvert(cls interface{}) Convert {
	s := Convert{}
	s.target = cls
	s.typ = reflect.TypeOf(cls)
	s.val = reflect.ValueOf(cls)
	return s
}

// func (sc *Convert) SetExcludes(s []string) {
// 	sc.excludes = s
// }

// SetTarget f
func (sc *Convert) SetTarget(cls interface{}) {
	sc.target = cls
	sc.typ = reflect.TypeOf(cls)
	sc.val = reflect.ValueOf(cls)
}

// GetTypeElem f
// func (sc *Convert) GetTypeElem() reflect.Type {
// 	if sc.typ.Kind() == reflect.Ptr {
// 		return sc.typ.Elem()
// 	}
// 	return sc.typ
// }

// GetType get reflect.Type
func (sc *Convert) GetType() reflect.Type {
	return sc.typ
}

// GetValue reflect.Value
func (sc *Convert) GetValue() reflect.Value {
	return sc.val
}

func (sc *Convert) buildFieldInfo() {
	typ := sc.typ
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	n := typ.NumField()

	sc.fields = make([]*StructField, 0)
	sc.dbfields = make([][]byte, 0)

	for i := 0; i < n; i++ {
		f := typ.Field(i)

		if f.Tag.Get("select") == "off" {
			continue
		}

		field := f.Tag.Get("json")

		if field == "" {
			if f.Name == skipFieldBaseVO || f.Name == skipFieldTable {
				continue
			}
			field = f.Name
		}

		sc.fields = append(sc.fields, &StructField{Name: f.Name, Field: field})
		sc.dbfields = append(sc.dbfields, []byte(field))
	}
}

// DBFields f
func (sc *Convert) DBFields() [][]byte {
	if len(sc.dbfields) == 0 {
		sc.buildFieldInfo()
	}
	return sc.dbfields
}

// Fields f
func (sc *Convert) Fields() []*StructField {
	if len(sc.fields) == 0 {
		sc.buildFieldInfo()
	}
	return sc.fields
}

func (sc *Convert) isEmpty(fieldTyp reflect.StructField, n int) bool {
	val := sc.val
	if val.Kind() == reflect.Ptr {
		val = sc.val.Elem()
	}
	return isEmptyValue(val.Field(n))
}

// ToMapRow to maprow
func (sc *Convert) ToMapRow(action int) MapRow {
	typ := sc.typ
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	l := typ.NumField()

	var skipTag, skipEmpty, tagtyp string
	var fieldVal reflect.Value
	isupdate := action == ActionUpdate

	actionType := ActionInsertStr
	if isupdate {
		actionType = ActionUpdateStr
	}

	maprow := MapRow{}

	for n := 0; n < l; n++ {
		field := typ.Field(n)
		dbField := field.Tag.Get(tagJSON)

		if dbField == "" {
			if field.Name == skipFieldBaseVO || field.Name == skipFieldTable {
				continue
			}
			dbField = field.Name
		}

		skipTag = field.Tag.Get("skip")
		if skipTag == "all" || skipTag == actionType {
			continue
		}

		if sc.val.Kind() == reflect.Ptr {
			fieldVal = sc.val.Elem().Field(n)
		} else {
			fieldVal = sc.val.Field(n)
		}
		value := fieldVal.Interface()

		skipEmpty = field.Tag.Get("skipempty")
		if (skipEmpty == "all" || skipEmpty == actionType) && (value == nil || sc.isEmpty(field, n)) {
			continue
		}

		// type
		tagtyp = field.Tag.Get(tagType)
		switch tagtyp {
		case typAutoTime:
			v := value.(time.Time)
			if v.IsZero() {
				maprow[dbField] = time.Now()
			} else {
				maprow[dbField] = v
			}
			continue

		case typPlus:
			if isupdate {
				maprow[dbField+"="+dbField+"+"+fmt.Sprint(value)] = nil
			}
			continue
		case typPlus1:
			if isupdate {
				maprow[dbField+"="+dbField+"+1"] = nil
			}
			continue
		}

		format := field.Tag.Get(tagFormat)
		if format != "" && field.Type.String() == timeType {
			value = (value.(time.Time)).Format(format)
		}

		maprow[dbField] = value
	}

	return maprow
}

// Row2VO f
func Row2VO(row MapRow, vo interface{}) {
	vof := reflect.ValueOf(vo)
	if vof.Kind() == reflect.Ptr {
		vof = vof.Elem()
	}
	typ := reflect.TypeOf(vo)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	n := typ.NumField()

	for i := 0; i < n; i++ {
		ftyp := typ.Field(i)
		fv := vof.Field(i)

		key := ftyp.Tag.Get("json")
		val, isok := row[key]
		if !isok {
			continue
		}
		// 如果同时存在数值

		switch ftyp.Type.Kind() {
		case reflect.Slice:

			switch ftyp.Type.String() {
			case "[]int":
				fv.Set(reflect.ValueOf(row.Ints(key)))
			case "[]float64":
				fv.Set(reflect.ValueOf(row.Floats(key)))
			case "[]int64":
				fv.Set(reflect.ValueOf(row.Int64s(key)))
			case "[]string":
				fv.Set(reflect.ValueOf(row.Strings(key)))
			default:
				fv.Set(reflect.ValueOf(val))
			}
			continue
		case reflect.Struct:
			switch ftyp.Type.String() {
			case timeType:
				if row.IsEmptyValue(key) {
					fv.Set(reflect.Zero(ftyp.Type))
					continue
				}
				// format is set
				format := ftyp.Tag.Get(tagFormat)
				if format == "" {
					fv.Set(reflect.ValueOf(row.Time(key)))
				} else {
					fv.Set(reflect.ValueOf(toTimeByFormat(row.String(key), format)))
				}
				continue

			default:
				prt := fv
				if fv.CanAddr() {
					prt = fv.Addr()
				}
				var err error
				switch row[key].(type) {
				case []byte:
					err = json.Unmarshal(row[key].([]byte), prt.Interface())
				case string:
					err = json.Unmarshal(util.Str2Bytes(row[key].(string)), prt.Interface())
				default:
					var src []byte
					src, _ = json.Marshal(row[key])
					err = json.Unmarshal(src, prt.Interface())
				}

				if err != nil {
					panic("Row2VO failed:" + err.Error() + "\nkey:" + key + "  type:" + reflect.TypeOf(row[key]).String())
				}
			}
			continue

		case reflect.Int, reflect.Int64, reflect.Int32:
			fv.SetInt(row.Int64(key))
			continue

		case reflect.String:
			switch row[key].(type) {
			case time.Time:
				t := row.Time(key)
				if t.IsZero() {
					fv.SetString("")
					continue
				}
				// format is set
				format := ftyp.Tag.Get("format")
				if format == "" {
					fv.SetString(t.Format(DateTimeFormat))
				} else {
					fv.SetString(t.Format(format))
				}

			default:
				fv.SetString(row.String(key))
			}
			continue

		default:
			fv.Set(reflect.ValueOf(val))
			continue
		}
		// fmt.Println(val, f.Type.String(), f.Type.Kind())
	}
}

// isEmptyValue check empty
func isEmptyValue(vof reflect.Value) bool {
	val := vof.Interface()
	switch val.(type) {
	case int:
		return val.(int) == 0
	case int64:
		return val.(int64) == 0
	case int32:
		return val.(int32) == 0
	case float32:
		return val.(float32) == 0
	case float64:
		return val.(float64) == 0
	case string:
		v := val.(string)
		return v == "" || v == "0"
	case time.Time:
		return (val.(time.Time)).IsZero()
	}

	switch vof.Type().Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return vof.Len() == 0
	}
	return false
}

func toTimeByFormat(src, format string) time.Time {
	if len(src) > len(format) {
		src = src[0:len(format)]
	}
	t, _ := time.Parse(format, src)
	return t
}

// StrictDBMapRow by VO, row from db query
func StrictDBMapRow(row MapRow, vo interface{}) {
	if row.IsEmpty() {
		return
	}

	vof := reflect.ValueOf(vo)
	if vof.Kind() == reflect.Ptr {
		vof = vof.Elem()
	}
	typ := reflect.TypeOf(vo)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	n := typ.NumField()

	for key := range row {
		ftyp, isok := typ.FieldByName(util.CamelCase(key))
		if !isok {
			for i := 0; i < n; i++ {
				ftyp = typ.Field(i)

				if ftyp.Tag.Get(tagJSON) == key {
					isok = true
					break
				}
			}
		}

		if !isok {
			delete(row, key)
			continue
		}
		if row.IsNull(key) {
			continue
		}

		switch ftyp.Type.Kind() {
		case reflect.Struct:
			switch ftyp.Type.String() {
			case timeType:
				if row.IsEmpty() {
					row[key] = ""
					continue
				}

				format := ftyp.Tag.Get("formart")
				if format == "" {
					row[key] = row.Time(key)
					continue
				}

				row[key] = toTimeByFormat(row.String(key), format)
			default:
				row[key] = util.Bytes2Str(row[key].([]byte))
			}

		case reflect.Slice:
			switch ftyp.Type.String() {
			case "[]int":
				row[key] = row.Ints(key)
			case "[]int64":
				row[key] = row.Int64s(key)
			case "[]float64":
				row[key] = row.Float(key)
			case "[]string":
				row[key] = row.Strings(key)
			default:
				panic("Strict db row error:" + ftyp.Type.String())
			}
		}
	}

}
