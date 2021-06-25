package dba

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/kere/gno/libs/util"
)

var (
	EmptyDataSet = DataSet{}
)

type DataSet struct {
	Fields  []string
	Columns [][]interface{}
	Types   []ColType
}

func NewDataSet(fields []string) DataSet {
	return DataSet{Fields: fields, Columns: make([][]interface{}, len(fields))}
}

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

func (d *DataSet) Len() int {
	if len(d.Columns) == 0 {
		return 0
	}
	return len(d.Columns[0])
}

func (d *DataSet) Release() {
	d.Fields = nil
	count := len(d.Columns)
	for i := 0; i < count; i++ {
		d.Columns[i] = nil
	}
	d.Columns = nil
	d.Types = nil
}

func (d *DataSet) AddRow(row []interface{}) {
	count := len(d.Columns)
	if count != len(row) {
		panic("db.AddRow columns.Len() != row.Len()")
	}
	for i := 0; i < count; i++ {
		d.Columns[i] = append(d.Columns[i], row[i])
	}
}

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

// ScanToDataSet db
func ScanToDataSet(rows *sql.Rows, isPool bool) (DataSet, error) {

	cols, err := rows.Columns()
	if err != nil {
		return EmptyDataSet, err
	}
	colsNum := len(cols)

	typs, err := rows.ColumnTypes()
	if err != nil {
		return EmptyDataSet, err
	}

	fields := make([]string, colsNum)
	typItems := make([]ColType, colsNum)
	for i := 0; i < colsNum; i++ {
		typItems[i] = NewColType(typs[i])
		fields[i] = typs[i].Name()
	}

	var result DataSet
	if isPool {
		result = GetDataSet(fields)
	} else {
		result.Fields = fields
		result.Columns = make([][]interface{}, colsNum)
	}
	result.Types = typItems

	// var row, tem []interface{}
	row := GetRow(colsNum)
	tem := GetRow(colsNum)
	defer PutRow(row)
	defer PutRow(tem)
	for i := 0; i < colsNum; i++ {
		tem[i] = &row[i]
	}

	for rows.Next() {
		if err = rows.Scan(tem...); err != nil {
			return result, err
		}
		result.AddRow(row)
	}

	return result, rows.Err()
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
