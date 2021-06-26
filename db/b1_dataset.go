package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/kere/gno/libs/util"
)

var (
	EmptyDataSet = DataSet{}
)

//DataSet
type DataSet struct {
	Fields  []string
	Columns [][]interface{}
	Types   []ColType
}

// NewDataSet
func NewDataSet(fields []string) DataSet {
	return DataSet{Fields: fields, Columns: make([][]interface{}, len(fields))}
}

// GetRow
func (d *DataSet) GetRow() []interface{} {
	return GetRow(len(d.Fields))
}

// RowAt
func (d *DataSet) RowAt(i int) []interface{} {
	n := d.Len()
	if n == 0 || i >= n {
		return nil
	}
	m := len(d.Columns)
	row := make([]interface{}, m)
	for k := 0; k < m; k++ {
		row[k] = d.Columns[k][i]
	}
	return row
}

// RowAtP
func (d *DataSet) RowAtP(i int) []interface{} {
	n := d.Len()
	if n == 0 || i >= n {
		return nil
	}
	m := len(d.Columns)
	row := GetRow(m)
	for k := 0; k < m; k++ {
		row[k] = d.Columns[k][i]
	}
	return row
}

// Len
func (d *DataSet) Len() int {
	if len(d.Columns) == 0 {
		return 0
	}
	return len(d.Columns[0])
}

// Release
func (d *DataSet) Release() {
	d.Fields = nil
	count := len(d.Columns)
	for i := 0; i < count; i++ {
		d.Columns[i] = nil
	}
	d.Columns = nil
	d.Types = nil
}

// AddRow
func (d *DataSet) AddRow(row []interface{}) {
	count := len(d.Columns)
	if count != len(row) {
		panic("db.AddRow columns.Len() != row.Len()")
	}
	for i := 0; i < count; i++ {
		d.Columns[i] = append(d.Columns[i], row[i])
	}
}

// AddRow0
func (d *DataSet) AddRow0(row []interface{}) {
	count := len(d.Columns)
	if count != len(row) {
		panic("db.AddRow columns.Len() != row.Len()")
	}
	for i := 0; i < count; i++ {
		d.Columns[i] = append(d.Columns[i], row[i])
	}
	PutRow(row)
}

// ColType class
type ColType struct {
	Name     string
	TypeName string
	Type     reflect.Type
	LengthOK bool
	Length   int

	DecimalOK bool
	Precision int
	Scale     int
}

// NewColType new type class
func NewColType(typ *sql.ColumnType) ColType {
	d := ColType{Name: typ.Name(), TypeName: typ.DatabaseTypeName(), Type: typ.ScanType()}
	var v, v2 int64
	v, d.LengthOK = typ.Length()
	d.Length = int(v)
	v, v2, d.DecimalOK = typ.DecimalSize()
	d.Precision = int(v)
	d.Scale = int(v2)
	return d
}

// PrintDataSet print
func PrintDataSet(dat *DataSet) {
	l := dat.Len()
	fmt.Println("------- length:", l, "-------")
	n := len(dat.Columns)
	for i := 0; i < n; i++ {
		fmt.Print(dat.Fields[i]+":"+dat.Types[i].TypeName, "\t")
	}
	fmt.Println()
	for i := 0; i < l; i++ {
		for k := 0; k < n; k++ {
			v := dat.Columns[k][i]
			switch v.(type) {
			case []byte:
				fmt.Print(util.Bytes2Str(v.([]byte)), "\t")
			default:
				fmt.Print(dat.Columns[k][i], "\t")
			}
		}
		fmt.Println()
	}
	fmt.Println("-- length:", l)
}

// DBRow
type DBRow struct {
	Values []interface{}
	Fields []string
}

// IsEmpty
func (d *DBRow) IsEmpty() bool {
	return len(d.Values) == 0
}

// Int
func (d *DBRow) Int(field string) int {
	v := d.Int64(field)
	return int(v)
}

// IntAt
func (d *DBRow) IntAt(i int) int {
	v := d.Int64At(i)
	return int(v)
}

// Int64At
func (d *DBRow) Int64(field string) int64 {
	i := util.StringsI(field, d.Fields)
	if i < 0 {
		return 0
	}
	return d.Int64At(i)
}

// Int64
func (d *DBRow) Int64At(i int) int64 {
	typ := reflect.TypeOf(d.Values[i])
	val := reflect.ValueOf(d.Values[i])
	switch typ.Kind() {
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return val.Int()
	case reflect.Slice:
		return 0
	case reflect.Float64, reflect.Float32:
		return int64(val.Float())
	case reflect.String:
		v, _ := strconv.ParseInt(val.String(), 10, 64)
		return v
	}
	return 0
}

// Float32
func (d *DBRow) Float32(field string) float32 {
	v := d.Float64(field)
	return float32(v)
}

// Float64
func (d *DBRow) Float64(field string) float64 {
	i := util.StringsI(field, d.Fields)
	if i < 0 {
		return 0
	}
	return d.Float64At(i)
}

// Float64At
func (d *DBRow) Float64At(i int) float64 {
	typ := reflect.TypeOf(d.Values[i])
	val := reflect.ValueOf(d.Values[i])
	switch typ.Kind() {
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return float64(val.Int())
	case reflect.Slice:
		return 0
	case reflect.Float64, reflect.Float32:
		return val.Float()
	case reflect.String:
		v, _ := strconv.ParseFloat(val.String(), 64)
		return v
	}
	return 0
}

// String
func (d *DBRow) String(field string) string {
	i := util.StringsI(field, d.Fields)
	if i < 0 {
		return ""
	}
	return d.StringAt(i)
}

// StringAt
func (d *DBRow) StringAt(i int) string {
	switch d.Values[i].(type) {
	case []byte:
		return util.Bytes2Str(d.Values[i].([]byte))
	case string:
		return d.Values[i].(string)
	}
	return fmt.Sprint(d.Values[i])
}

// Bytes
func (d *DBRow) Bytes(field string) []byte {
	i := util.StringsI(field, d.Fields)
	if i < 0 {
		return nil
	}
	return d.BytesAt(i)
}

// BytesAt
func (d *DBRow) BytesAt(i int) []byte {
	switch d.Values[i].(type) {
	case []byte:
		return d.Values[i].([]byte)
	case string:
		return util.Str2Bytes(d.Values[i].(string))
	}
	return util.Str2Bytes(fmt.Sprint(d.Values[i]))
}
