package db

import (
	"fmt"
	"reflect"
	"time"
)

// StructField class
type StructField struct {
	Field string
	Name  string
}

// StructConverter class
// Example: sc := NewStructConvert(UserVO{})
// dataset, err := db.NewQueryBuilder("users").Where("id=?",1).Query()
// result := cs.DataSet2Struct(maprow)
// userVO := result[0].(*UserVO)
type StructConverter struct {
	typ      reflect.Type
	val      reflect.Value
	target   interface{}
	fields   []*StructField
	dbfields [][]byte
}

// NewStructConvert func
func NewStructConvert(cls interface{}) *StructConverter {
	s := &StructConverter{}
	s.SetTarget(cls)
	return s
}

// func (sc *StructConverter) SetExcludes(s []string) {
// 	sc.excludes = s
// }

// SetTarget f
func (sc *StructConverter) SetTarget(cls interface{}) {
	sc.target = cls
	sc.typ = reflect.TypeOf(cls)

	sc.val = reflect.ValueOf(cls)
	// if sc.val.Kind() == reflect.Ptr {
	// 	sc.val = sc.val.Elem()
	// }
}

// GetTypeElem f
func (sc *StructConverter) GetTypeElem() reflect.Type {
	if sc.typ.Kind() == reflect.Ptr {
		return sc.typ.Elem()
	}
	return sc.typ
}

// GetType get reflect.Type
func (sc *StructConverter) GetType() reflect.Type {
	return sc.typ
}

// GetValue reflect.Value
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

// DBFields f
func (sc *StructConverter) DBFields() [][]byte {
	if len(sc.dbfields) == 0 {
		sc.buildFieldInfo()
	}
	return sc.dbfields
}

// Fields f
func (sc *StructConverter) Fields() []*StructField {
	if len(sc.fields) == 0 {
		sc.buildFieldInfo()
	}
	return sc.fields
}

// // KeyValueList f
// func (sc *StructConverter) KeyValueList(stype string) ([][]byte, []interface{}, [][]byte) {
// 	return keyValueList(stype, sc.target)
// }

//DataRow2Struct f
func (sc *StructConverter) DataRow2Struct(maprow MapRow) (IVO, error) {
	rowStruct := reflect.New(sc.GetTypeElem())
	err := maprow.CopyToVO(rowStruct.Interface().(IVO))
	if err != nil {
		return nil, err
	}
	return (rowStruct.Interface()).(IVO), nil
}

// // DataSet2Struct f
// func (sc *StructConverter) DataSet2Struct(dataset DataSet) (VODataSet, error) {
// 	result := VODataSet{}
// 	var err error
// 	var vo IVO
//
// 	for _, row := range dataset {
// 		vo, err = sc.DataRow2Struct(row)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result = append(result, vo)
// 	}
//
// 	return result, nil
// }

func (sc *StructConverter) isEmpty(fieldTyp reflect.StructField, n int) bool {
	val := sc.val
	if val.Kind() == reflect.Ptr {
		val = sc.val.Elem()
	}

	switch fieldTyp.Type.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		if val.Field(n).Len() == 0 {
			return true
		}
	case reflect.Struct:
		if fieldTyp.Type.String() == timeClassName && (val.Field(n).Interface().(time.Time)).IsZero() {
			return true
		}

		return !val.Field(n).IsValid()

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
		if s == vNil || s == "" || s == "0" {
			return true
		}
	}
	return false
}

const (
	tagSkip      = "skip"
	tagJSON      = "json"
	tagSkipEmpty = "skipempty"
	tagType      = "type"
	typTime      = "time"
	typAutoTime  = "autotime"
	typPlus      = "plus"
	typPlus1     = "plus1"
	typVersion   = "version"
	vNil         = "<nil>"
	vAll         = "all"
	vBaseVO      = "BaseVO"
)

// Struct2DataRow to maprow
func (sc *StructConverter) Struct2DataRow(action int) MapRow {
	typ := sc.GetTypeElem()
	l := typ.NumField()

	var skipTag, skipEmpty, tagtyp string
	var value interface{}
	var fieldVal reflect.Value
	isupdate := action == ActionUpdate

	actionType := "insert"
	if isupdate {
		actionType = "update"
	}

	maprow := MapRow{}

	for n := 0; n < l; n++ {
		field := typ.Field(n)
		dbField := field.Tag.Get(tagJSON)

		if dbField == "" {
			if field.Name == vBaseVO {
				continue
			}
			dbField = field.Name
		}

		skipTag = field.Tag.Get(tagSkip)
		if skipTag == vAll || skipTag == actionType {
			continue
		}

		if sc.val.Kind() == reflect.Ptr {
			fieldVal = sc.val.Elem().Field(n)
		} else {
			fieldVal = sc.val.Field(n)
		}
		value = fieldVal.Interface()

		tagtyp = field.Tag.Get(tagType)
		switch tagtyp {
		case typTime:
			v := value.(time.Time)
			if v.IsZero() {
				value = nil
			}
		case typAutoTime:
			v := value.(time.Time)
			if v.IsZero() {
				value = time.Now()
			}
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

		if fieldVal.Kind() == reflect.Ptr && fmt.Sprint(fieldVal) == vNil {
			value = nil
		}

		skipEmpty = field.Tag.Get(tagSkipEmpty)
		if (skipEmpty == "all" || skipEmpty == actionType) && (value == nil || sc.isEmpty(field, n)) {
			continue
		}

		maprow[dbField] = value
	}

	return maprow
}
