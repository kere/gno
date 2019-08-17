package db

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// MapRows datarow list
type MapRows []MapRow

// NewMapRows by length
func NewMapRows(l int) MapRows {
	return make([]MapRow, l)
}

// NewMapRowsN by length
func NewMapRowsN(l, n int) MapRows {
	return make([]MapRow, l, n)
}

// Insert current dataset
func (ds MapRows) Insert(table string) error {
	if ds.Len() == 0 {
		return nil
	}
	ins := InsertBuilder{}
	_, err := ins.Table(table).InsertM(ds)
	return err
}

// IsEmpty func
func (ds MapRows) IsEmpty() bool {
	return len(ds) == 0
}

// Clone func
func (ds MapRows) Clone() MapRows {
	tmp := make([]MapRow, len(ds))
	for i := range ds {
		tmp[i] = ds[i].Clone()
	}
	return tmp
}

// Bytes2String fix bytes value data
func (ds MapRows) Bytes2String() MapRows {
	for i := range ds {
		ds[i] = ds[i].Bytes2String()
	}

	return ds
}

// // Bytes2NumericBySubfix fix number
// // Convert bytes data to numeric by subfix
// func (ds MapRows) Bytes2NumericBySubfix(subfix string) MapRows {
// 	for i := range ds {
// 		ds[i] = ds[i].Bytes2NumericBySubfix(subfix)
// 	}
//
// 	return ds
// }

// // Bytes2NumericByFields Convert bytes data to numeric by field list
// func (ds MapRows) Bytes2NumericByFields(fields []string) MapRows {
// 	for i := range ds {
// 		ds[i] = ds[i].Bytes2NumericByFields(fields)
// 	}
//
// 	return ds
// }

// Search search by field and return MapRow, if not found return nil
func (ds MapRows) Search(field string, value interface{}) MapRow {
	for _, r := range ds {
		if r[field] == value {
			return r
		}
	}
	return nil
}

// Filter data
func (ds MapRows) Filter(f func(int, MapRow) bool) MapRows {
	arr := MapRows{}
	for i, r := range ds {
		if f(i, r) {
			arr = append(arr, r)
		}
	}
	return arr
}

// Group data
func (ds MapRows) Group(field string) map[interface{}]MapRows {
	mapDat := make(map[interface{}]MapRows, 0)
	for _, r := range ds {
		v := r[field]
		if _, isok := mapDat[v]; !isok {
			mapDat[v] = MapRows{}
		}

		mapDat[v] = append(mapDat[v], r)
	}
	return mapDat
}

// Len return length of dataset
func (ds MapRows) Len() int {
	return len(ds)
}

// StringValues return string value list by field
func (ds MapRows) StringValues(field string) []string {
	l := ds.Len()
	arr := make([]string, l)
	for i := 0; i < l; i++ {
		arr[i] = ds[i].String(field)
	}
	return arr
}

// FloatValues return string value list by field
func (ds MapRows) FloatValues(field string) []float64 {
	l := ds.Len()
	arr := make([]float64, l)
	for i := 0; i < l; i++ {
		arr[i] = ds[i].Float(field)
	}
	return arr
}

// IntValues return string value list by field
func (ds MapRows) IntValues(field string) []int {
	l := ds.Len()
	arr := make([]int, l)
	for i := 0; i < l; i++ {
		arr[i] = ds[i].Int(field)
	}
	return arr
}

// DataSet zip dataset
func (ds MapRows) DataSet(fields []string) DataSet {
	var dataset DataSet
	l := len(ds)
	if l == 0 {
		return dataset
	}
	// dataset = DataSet{Fields: fields}
	n := len(ds[0])
	cols := NewColumns(n, l)
	for i := 0; i < l; i++ {
		cols.SetMapRow(i, ds[i], fields)
	}
	dataset = DataSet{Fields: fields, Columns: cols}

	return dataset
}

// VORows class
type VORows []IVO

// ToJSON []byte
func (ds VORows) ToJSON(action int) []byte {
	l := len(ds)
	buf := bytes.NewBuffer([]byte("["))

	for i := 0; i < l; i++ {
		if ds[i] == nil {
			continue
		}

		row := ds[i].ToMapRow(action)
		src, err := json.Marshal(row)
		if err != nil {
			continue
		}
		if i > 0 {
			buf.Write(BCommaSplit)
		}
		buf.Write(src)
	}

	buf.WriteString("]")

	return buf.Bytes()
}

// Encode to
func (ds VORows) Encode() [][]interface{} {
	count := len(ds)
	if count < 1 {
		return make([][]interface{}, 0)
	}

	sc := NewStructConvert(ds[0])
	dbFields := sc.DBFields()
	fields := sc.Fields()
	colSize := len(dbFields)
	values := make([][]interface{}, len(ds)+1)
	columns := make([]interface{}, colSize)

	var i int
	var k int
	var tmp []interface{}

	for i = range dbFields {
		columns[i] = string(dbFields[i])
	}
	values[0] = columns

	for i = 0; i < count; i++ {
		tmp = make([]interface{}, colSize)
		for k = range fields {
			tmp[k] = reflect.ValueOf(ds[i]).Elem().FieldByName(fields[k].Name).Interface()
		}

		values[i+1] = tmp
	}

	return values
}
