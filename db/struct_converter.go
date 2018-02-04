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
	excludes []string
	fields   []*StructField
	dbfields [][]byte
}

func NewStructConvert(cls interface{}) *StructConverter {
	s := &StructConverter{}
	s.SetTarget(cls)
	return s
}

func (sc *StructConverter) SetExcludes(s []string) {
	sc.excludes = s
}

func (sc *StructConverter) SetTarget(cls interface{}) {
	sc.target = cls
	sc.typ = reflect.TypeOf(cls)

	sc.val = reflect.ValueOf(cls)
	if sc.val.Kind() == reflect.Ptr {
		sc.val = sc.val.Elem()
	}
}

func (sc *StructConverter) GetTypeElem() reflect.Type {
	if sc.typ.Kind() == reflect.Ptr {
		return sc.typ.Elem()
	} else {
		return sc.typ
	}
}

func (sc *StructConverter) GetType() reflect.Type {
	return sc.typ
}

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

		if InStrings(sc.excludes, field) {
			continue
		}

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
	return keyValueList(stype, sc.target, sc.excludes)
}

func (sc *StructConverter) DataRow2Struct(datarow DataRow) (IVO, error) {
	rowStruct := reflect.New(sc.GetTypeElem())
	err := datarow.ConvertTo(rowStruct.Interface().(IVO))
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
	switch fieldTyp.Type.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		if sc.val.Field(n).Len() == 0 {
			return true
		}
	case reflect.Struct:
		return !sc.val.Field(n).IsValid()

	case reflect.Interface:
		return sc.val.Field(n).IsNil()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return sc.val.Field(n).Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return sc.val.Field(n).Uint() == 0

	case reflect.Float32, reflect.Float64:
		return sc.val.Field(n).Float() == 0

	case reflect.String:
		return sc.val.Field(n).String() == ""

	default:
		s := fmt.Sprint(sc.val.Field(n).Interface())
		if s == "" || s == "0" {
			return true
		}
	}
	return false
}

func (sc *StructConverter) Struct2DataRow(actionType string) DataRow {
	typ := sc.GetTypeElem()

	l := typ.NumField()

	var field reflect.StructField
	var dbField string

	skipTag := ""
	skipEmpty := ""
	autotime := ""
	// isJSONType := false
	var value interface{}

	datarow := DataRow{}

	for n := 0; n < l; n++ {
		field = typ.Field(n)
		dbField = field.Tag.Get("json")

		if dbField == "" {
			if field.Name == "BaseVO" {
				continue
			}
			dbField = field.Name
		}

		// isJSONType = len(dbField) > 5 && dbField[len(dbField)-5:] == "_json"

		skipTag = field.Tag.Get("skip")
		if actionType != "" && skipTag == actionType {
			continue
		}

		if InStrings(sc.excludes, dbField) || skipTag == "all" {
			continue
		}

		value = sc.val.Field(n).Interface()

		skipEmpty = field.Tag.Get("skipempty")
		autotime = field.Tag.Get("autotime")
		if autotime == "true" && field.Type.String() == "time.Time" && (value.(time.Time)).IsZero() {
			if Current().Location.String() == "UTC" {
				datarow[dbField] = time.Now()
			} else {
				datarow[dbField] = time.Now().In(Current().Location)
			}
			continue
		}

		if (skipEmpty == "all" || skipEmpty == actionType) && sc.isEmpty(field, n) {
			continue
		}

		datarow[dbField] = value

		// else if isJSONType {
		// 	b, err := json.Marshal(value)
		// 	if err != nil {
		// 		Current().Log.Fatalln("Json Marshal ", err)
		// 		b = []byte("")
		// 	}
		//
		// 	if Current().Driver.DriverName() == "pgsql" {
		// 		datarow[dbField] = b
		// 	} else {
		// 		datarow[dbField] = string(b)
		// 	}
		//
		// }
	}

	return datarow
}
