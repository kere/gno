package db

import "reflect"

// DataSet datarow list
type DataSet []DataRow

// Insert current dataset
func (ds DataSet) Insert(table string) error {
	if ds.Len() == 0 {
		return nil
	}
	_, err := NewInsertBuilder(table).InsertM(ds)
	return err
}

// IsEmpty func
func (ds DataSet) IsEmpty() bool {
	return len(ds) == 0
}

// Clone func
func (ds DataSet) Clone() DataSet {
	tmp := make([]DataRow, len(ds))
	for i := range ds {
		tmp[i] = ds[i].Clone()
	}
	return tmp
}

// Bytes2String fix bytes value data
func (ds DataSet) Bytes2String() DataSet {
	for i := range ds {
		ds[i] = ds[i].Bytes2String()
	}

	return ds
}

// Bytes2NumericBySubfix fix number
// Convert bytes data to numeric by subfix
func (ds DataSet) Bytes2NumericBySubfix(subfix string) DataSet {
	for i := range ds {
		ds[i] = ds[i].Bytes2NumericBySubfix(subfix)
	}

	return ds
}

// Bytes2NumericByFields Convert bytes data to numeric by field list
func (ds DataSet) Bytes2NumericByFields(fields []string) DataSet {
	for i := range ds {
		ds[i] = ds[i].Bytes2NumericByFields(fields)
	}

	return ds
}

// Search search by field and return DataRow, if not found return nil
func (ds DataSet) Search(field string, value interface{}) DataRow {
	for _, r := range ds {
		if r[field] == value {
			return r
		}
	}
	return nil
}

// Filter data
func (ds DataSet) Filter(f func(int, DataRow) bool) DataSet {
	arr := DataSet{}
	for i, r := range ds {
		if f(i, r) {
			arr = append(arr, r)
		}
	}
	return arr
}

// Group data
func (ds DataSet) Group(field string) map[interface{}]DataSet {
	mapDat := make(map[interface{}]DataSet, 0)
	for _, r := range ds {
		v := r[field]
		if _, isok := mapDat[v]; !isok {
			mapDat[v] = DataSet{}
		}

		mapDat[v] = append(mapDat[v], r)
	}
	return mapDat
}

// Len return length of dataset
func (ds DataSet) Len() int {
	return len(ds)
}

// StringValues return string value list by field
func (ds DataSet) StringValues(field string) []string {
	l := ds.Len()
	arr := make([]string, l)
	for i := 0; i < l; i++ {
		arr[i] = ds[i].String(field)
	}
	return arr
}

// FloatValues return string value list by field
func (ds DataSet) FloatValues(field string) []float64 {
	l := ds.Len()
	arr := make([]float64, l)
	for i := 0; i < l; i++ {
		arr[i] = ds[i].Float(field)
	}
	return arr
}

// IntValues return string value list by field
func (ds DataSet) IntValues(field string) []int {
	l := ds.Len()
	arr := make([]int, l)
	for i := 0; i < l; i++ {
		arr[i] = ds[i].Int(field)
	}
	return arr
}

// Encode zip dataset
func (ds DataSet) Encode() [][]interface{} {
	ds.Bytes2String()
	if ds == nil || len(ds) == 0 {
		return nil
	}

	var value []interface{}

	values := make([][]interface{}, len(ds)+1)
	colSize := len(ds[0])
	columns := make([]interface{}, colSize)
	colMap := make(map[string]int)
	isFirst := true
	var n int
	var row DataRow
	for i, r := range ds {
		row = r
		n = 0
		if isFirst {
			for k := range row {
				columns[n] = k
				colMap[k] = n
				n++
			}

			values[i] = columns
			isFirst = false
		}

		value = make([]interface{}, colSize)
		for k, v := range row {
			value[colMap[k]] = v
		}

		values[i+1] = value
	}

	return values
}

type VODataSet []IVO

func (ds VODataSet) Encode() [][]interface{} {
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
