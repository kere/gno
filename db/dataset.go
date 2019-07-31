package db

import (
	"database/sql"
	"errors"
	"reflect"
)

// DataType class
type DataType struct {
	Name     string
	TypeName string
	Type     reflect.Type
	LengthOK bool
	Length   int

	DecimalOK bool
	Precision int
	Scale     int
}

// NewDataType new type class
func NewDataType(typ *sql.ColumnType) DataType {
	d := DataType{Name: typ.Name(), TypeName: typ.DatabaseTypeName(), Type: typ.ScanType()}
	var v, v2 int64
	v, d.LengthOK = typ.Length()
	d.Length = int(v)
	v, v2, d.DecimalOK = typ.DecimalSize()
	d.Precision = int(v)
	d.Scale = int(v2)
	return d
}

// DataRow class
type DataRow []interface{}

// DataRow2MapRow convert to
func DataRow2MapRow(row DataRow, fields []string) MapRow {
	mapRow := MapRow{}
	n := len(fields)
	for i := 0; i < n; i++ {
		mapRow[fields[i]] = row[i]
	}
	return mapRow
}

// DataColumn class
type DataColumn []interface{}

// DataSet data rows
type DataSet struct {
	Fields  []string
	Types   []DataType
	Columns []DataColumn
}

// Len dataset
func (d *DataSet) Len() int {
	if len(d.Columns) == 0 {
		return -1
	}
	return len(d.Columns[0])
}

// First datarow
func (d *DataSet) First() MapRow {
	return d.MapRowAt(0)
}

// MapRowAt datarow
func (d *DataSet) MapRowAt(i int) MapRow {
	n := len(d.Columns)
	if n == 0 || len(d.Columns[0]) >= i {
		return nil
	}

	row := MapRow{}

	for k := 0; k < n; k++ {
		row[d.Fields[k]] = d.Columns[k][i]
	}
	return row
}

// AddMapRow add
func (d *DataSet) AddMapRow(row MapRow) error {
	n := len(d.Types)
	if n == 0 {
		return errors.New("dataset types is empty")
	}
	if n != len(row) {
		return errors.New("row fields not match")
	}

	if len(d.Columns) == 0 {
		d.Columns = make([]DataColumn, n)
	}

	for k := 0; k < n; k++ {
		val, isok := row[d.Fields[k]]
		if isok {
			d.Columns[k] = append(d.Columns[k], val)
		}
	}

	return nil
}

// AddDataRow add
func (d *DataSet) AddDataRow(row DataRow) error {
	n := len(row)
	if n == 0 {
		return nil
	}
	if n != len(d.Fields) {
		return errors.New("dataset AddDataRow, fields not matched")
	}
	for k := 0; k < n; k++ {
		d.Columns[k] = append(d.Columns[k], row[k])
	}
	return nil
}
