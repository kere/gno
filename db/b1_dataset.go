package db

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	iconv "github.com/djimenez/iconv-go"
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
func (d *DataSet) GetRowP() []interface{} {
	return GetRow(len(d.Fields))
}

// GetDBRow
func (d *DataSet) GetDBRow() DBRow {
	return DBRow{Fields: d.Fields, Values: GetRow(len(d.Fields))}
}

// RowAt
func (d *DataSet) RowAt(i int, row []interface{}) {
	n := d.Len()
	if n == 0 || i >= n {
		return
	}
	m := len(d.Columns)
	if m != len(row) {
		panic("DataSet RowAt: columns.Len() != row.Len()")
	}
	for k := 0; k < m; k++ {
		row[k] = d.Columns[k][i]
	}
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

// DBRowAt
func (d *DataSet) DBRowAt(i int, dbRow DBRow) {
	n := d.Len()
	if n == 0 || i >= n {
		return
	}
	m := len(d.Columns)
	if m != len(dbRow.Values) {
		panic("DataSet RowAt: columns.Len() != row.Len()")
	}
	for k := 0; k < m; k++ {
		dbRow.Values[k] = d.Columns[k][i]
	}
}

// DBRowAtP
func (d *DataSet) DBRowAtP(i int) DBRow {
	n := d.Len()
	if n == 0 || i >= n {
		return DBRow{}
	}

	m := len(d.Fields)
	vals := GetRow(m)
	r := DBRow{Fields: d.Fields, Values: vals}
	for k := 0; k < m; k++ {
		vals[k] = d.Columns[k][i]
	}
	return r
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
		row[i] = nil
	}
}

// SetRow
func (d *DataSet) SetRow(i int, row []interface{}) {
	count := len(d.Columns)
	if count != len(row) {
		panic("db.SetRow columns.Len() != row.Len()")
	}
	for k := 0; k < count; k++ {
		d.Columns[k][i] = row[k]
	}
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

// PrintRow print
func PrintRow(r []interface{}) {
	l := len(r)
	fmt.Printf("-- %d :", l)
	for i := 0; i < l; i++ {
		v := r[i]
		switch v.(type) {
		case []byte:
			fmt.Print(util.Bytes2Str(v.([]byte)), util.STab)
		default:
			fmt.Print(r[i], util.STab)
		}
	}
	fmt.Println()
}

// PrintDataSet print
func PrintDataSet(ds *DataSet) {
	l := ds.Len()
	fmt.Println("------- length:", l, "-------")
	n := len(ds.Columns)
	spl := util.STab
	if len(ds.Types) == 0 {
		for i := 0; i < n; i++ {
			fmt.Print(ds.Fields[i] + spl)
		}
	} else {
		for i := 0; i < n; i++ {
			fmt.Print(ds.Fields[i]+":"+ds.Types[i].TypeName, spl)
		}
	}
	fmt.Println()
	for i := 0; i < l; i++ {
		for k := 0; k < n; k++ {
			v := ds.Columns[k][i]
			switch v.(type) {
			case []byte:
				fmt.Print(util.Bytes2Str(v.([]byte)), spl)
			default:
				fmt.Print(ds.Columns[k][i], spl)
			}
		}
		fmt.Println()
	}
	fmt.Println("-- length:", l)
}

// RangeI
func (d *DataSet) IRange(a, b int) DataSet {
	l := d.Len()
	if l == 0 {
		return *d
	}

	if b == -1 || b > l-1 {
		b = l - 1
	}
	if a == -1 || a > b {
		return EmptyDataSet
	}

	ds := NewDataSet(d.Fields)
	n := len(ds.Columns)
	for i := 0; i < n; i++ {
		ds.Columns[i] = d.Columns[i][a : b+1]
	}
	return ds
}

// EachPage
func (d *DataSet) EachPage(pageSize int, f func(page int, ds DataSet) bool) int {
	count := d.Len()
	ll := len(d.Columns)
	cols := make([][]interface{}, ll)
	n := util.EachPage(count, pageSize, func(pageN int, a int, b int) bool {
		ds := DataSet{Fields: d.Fields, Columns: cols}
		for i := 0; i < ll; i++ {
			cols[i] = d.Columns[i][a:b]
		}
		ds.Columns = cols
		return f(pageN, ds)
	})
	return n
}

func LoadCSV(filename string, hasFields bool) (DataSet, error) {
	ds := DataSet{}
	var err error
	ds.Fields, ds.Columns, err = loadCSV(filename, hasFields, false)
	return ds, err
}

func LoadCSVP(filename string, hasFields bool) (DataSet, error) {
	ds := DataSet{}
	var err error
	ds.Fields, ds.Columns, err = loadCSV(filename, hasFields, true)
	return ds, err
}

func loadCSV(filename string, hasFields, isPool bool) ([]string, [][]interface{}, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	sep := util.BComma
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	src := scanner.Bytes()
	isGB := util.IsGBK(src)
	toSrc := make([]byte, len(src)*2)
	var wn int
	if isGB {
		_, wn, _ = iconv.Convert(src, toSrc, "gb2312", "utf-8")
		src = toSrc[:wn]
	}

	var columns [][]interface{}
	var fields []string

	var ll int
	if hasFields {
		fields = strings.Split(string(src), util.SComma)
		// fields = util.SplitStrNotSafe(util.Bytes2Str(src), util.SComma)
		ll = len(fields)
		columns = make([][]interface{}, ll)
		if isPool {
			for i := 0; i < ll; i++ {
				columns[i] = GetColumn()
			}
		}
	} else {
		arr := util.SplitBytesNotSafe(src, sep)
		// arr := bytes.Split(src, sep)
		ll = len(arr)
		fields = make([]string, ll)
		columns = make([][]interface{}, ll)
		for i := 0; i < ll; i++ {
			if isPool {
				columns[i] = GetColumn()
			}
			fields[i] = fmt.Sprint("val", i+1)
			columns[i] = append(columns[i], string(bytes.TrimSpace(arr[i])))
			// columns[i] = append(columns[i], bytes.TrimSpace(arr[i]))
		}
	}

	for scanner.Scan() {
		src := scanner.Bytes()
		if len(src) == 0 {
			continue
		}
		if isGB {
			for i := 0; i < wn; i++ {
				toSrc[i] = 0
			}
			_, wn, _ := iconv.Convert(src, toSrc, "gb2312", "utf-8")
			src = toSrc[:wn]
		}
		// fmt.Println(string(src))
		// arr := bytes.Split(src, sep)
		arr := util.SplitBytesNotSafe(src, sep)
		if len(arr) != ll {
			continue
		}

		for k := 0; k < ll; k++ {
			columns[k] = append(columns[k], string(bytes.TrimSpace(arr[k])))
		}
	}

	return fields, columns, nil
}
