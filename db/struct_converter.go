package db

import (
	"fmt"
	"reflect"
	"time"
)

type StructField struct {
	Field string
	Name  string
}

// Example: sc := NewStructConvert(UserVO{})
// dataset, err := db.NewQueryBuilder("users").Where("id=?",1).Query()
// result := cs.DataSet2Struct(datarow)
// userVO := result[0].(*UserVO)
type StructConverter struct {
	typ      reflect.Type
	val      reflect.Value
	target   interface{}
	fields   []*StructField
	dbfields [][]byte
}

func NewStructConvert(cls interface{}) *StructConverter {
	s := &StructConverter{}
	s.SetTarget(cls)
	return s
}

// func (sc *StructConverter) SetExcludes(s []string) {
// 	sc.excludes = s
// }

func (sc *StructConverter) SetTarget(cls interface{}) {
	sc.target = cls
	sc.typ = reflect.TypeOf(cls)

	sc.val = reflect.ValueOf(cls)
	// if sc.val.Kind() == reflect.Ptr {
	// 	sc.val = sc.val.Elem()
	// }
}

func (sc *StructConverter) GetTypeElem() reflect.Type {
	if sc.typ.Kind() == reflect.Ptr {
		return sc.typ.Elem()
	} else {
		return sc.typ
	}
}

// GetType reflect.Type
func (sc *StructConverter) GetType() reflect.Type {
	return sc.typ
}

// GetType reflect.Value
func (sc *StructConverter) GetValue() reflect.Value {
	return sc.val
}

func (sc *StructConverter) buildFieldInfo() {
	typ := sc.GetTypeElem()
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
			if f.Name == "BaseVO" {
				continue
			}
			field = f.Name
		}

		// if InStrings(sc.excludes, field) {
		// 	continue
		// }

		sc.fields = append(sc.fields, &StructField{Name: f.Name, Field: field})
		sc.dbfields = append(sc.dbfields, []byte(field))
	}
}

func (sc *StructConverter) DBFields() [][]byte {
	if len(sc.dbfields) == 0 {
		sc.buildFieldInfo()
	}
	return sc.dbfields
}

func (sc *StructConverter) Fields() []*StructField {
	if len(sc.fields) == 0 {
		sc.buildFieldInfo()
	}
	return sc.fields
}

func (sc *StructConverter) KeyValueList(stype string) ([][]byte, []interface{}, [][]byte) {
	return keyValueList(stype, sc.target)
}

func (sc *StructConverter) DataRow2Struct(datarow DataRow) (IVO, error) {
	rowStruct := reflect.New(sc.GetTypeElem())
	err := datarow.CopyToVO(rowStruct.Interface().(IVO))
	if err != nil {
		return nil, err
	}
	return (rowStruct.Interface()).(IVO), nil
}

func (sc *StructConverter) DataSet2Struct(dataset DataSet) (VODataSet, error) {
	result := VODataSet{}
	var err error
	var vo IVO

	for _, row := range dataset {
		vo, err = sc.DataRow2Struct(row)
		if err != nil {
			return nil, err
		}
		result = append(result, vo)
	}

	return result, nil
}

func (sc *StructConverter) isEmpty(fieldTyp reflect.StructField, n int) bool {
	val := sc.val
	if val.Kind() == reflect.Ptr {
		val = sc.val.Elem()
	}

	switch fieldTyp.Type.Kind() {
	case reflect.Ptr:
		return val.Field(n).IsNil()

	case reflect.Map, reflect.Slice, reflect.Array:
		if val.Field(n).Len() == 0 {
			return true
		}
	case reflect.Struct:
		if fieldTyp.Type.String() == timeClassName && (val.Field(n).Interface().(time.Time)).IsZero() {
			return true
		}

		return !val.Field(n).IsValid()

	case reflect.Interface:
		return val.Field(n).IsNil()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Field(n).Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Field(n).Uint() == 0

	case reflect.Float32, reflect.Float64:
		return val.Field(n).Float() == 0

	case reflect.String:
		return val.Field(n).String() == ""

	default:
		s := fmt.Sprint(val.Field(n).Interface())
		if s == "" || s == "0" {
			return true
		}
	}
	return false
}

const (
	tagSkip      = "skip"
	tagJson      = "json"
	tagPlus      = "plus"
	tagCType     = "ctype"
	tagSkipEmpty = "skipempty"
	tagAutotime  = "autotime"
	vTrue        = "true"
	vAll         = "all"
	vBaseVO      = "BaseVO"
)

// Struct2DataRow to datarow
func (sc *StructConverter) Struct2DataRow(actionType string) DataRow {
	typ := sc.GetTypeElem()

	l := typ.NumField()

	var skipTag, skipEmpty, autotime string
	var value interface{}

	datarow := DataRow{}

	for n := 0; n < l; n++ {
		field := typ.Field(n)
		dbField := field.Tag.Get(tagJson)

		if dbField == "" {
			if field.Name == vBaseVO {
				continue
			}
			dbField = field.Name
		}

		skipTag = field.Tag.Get(tagSkip)
		if actionType != "" && skipTag == actionType {
			continue
		}

		if skipTag == vAll {
			continue
		}

		if sc.val.Kind() == reflect.Ptr {
			value = sc.val.Elem().Field(n).Interface()
		} else {
			value = sc.val.Field(n).Interface()
		}

		// plus field
		tagplus := field.Tag.Get(tagPlus)
		if actionType == ActionUpdate && tagplus != "" {
			// version=version+1
			datarow[dbField+"="+dbField+"+"+tagplus] = nil
			continue
		}

		skipEmpty = field.Tag.Get(tagSkipEmpty)
		autotime = field.Tag.Get(tagAutotime)
		if (autotime == vTrue || autotime == actionType) && field.Type.String() == "time.Time" && (value.(time.Time)).IsZero() {
			// ------- time zone --------
			datarow[dbField] = time.Now()
			continue
		}

		if (skipEmpty == "" || skipEmpty == "all" || skipEmpty == actionType) && sc.isEmpty(field, n) {
			continue
		}

		datarow[dbField] = value
	}

	return datarow
}
