package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/kere/gno/libs/util"
)

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

// Row class
type Row []interface{}

// Row2MapRow convert to
func Row2MapRow(row Row, fields []string) MapRow {
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
	Table       string         `json:"table"`
	Fields      []string       `json:"fields"`
	Types       []ColType      `json:"-"`
	Columns     []DataColumn   `json:"columns"`
	fieldIndexs map[string]int // 缓存field 索引
}

// PrintDataSet print
func PrintDataSet(dat *DataSet) {
	l := dat.Len()
	n := len(dat.Columns)
	fmt.Println(strings.Join(dat.Fields, "\t"))
	for i := 0; i < l; i++ {
		for k := 0; k < n; k++ {
			v := dat.Columns[k][i]
			switch v.(type) {
			case []byte:
				fmt.Print(util.BytesToStr(v.([]byte)), "\t")
			default:
				fmt.Print(dat.Columns[k][i], "\t")
			}
		}
		fmt.Println()
	}
	fmt.Println("length:", l)
}

// NewDataSet new
func NewDataSet(fields []string) DataSet {
	return DataSet{Fields: fields, Columns: make([]DataColumn, len(fields))}
}

// Len dataset
func (d *DataSet) Len() int {
	if len(d.Columns) == 0 {
		return -1
	}
	return len(d.Columns[0])
}

// FieldI 字段索引
func (d *DataSet) FieldI(field string) int {
	if d.fieldIndexs == nil {
		d.fieldIndexs = make(map[string]int)
	}
	n := len(d.Fields)
	for i := 0; i < n; i++ {
		if d.Fields[i] == field {
			d.fieldIndexs[field] = i
			return i
		}
	}
	return -1
}

// First datarow
func (d *DataSet) First() MapRow {
	return d.MapRowAt(0)
}

// RangeI datarow
func (d *DataSet) RangeI(b, e int) DataSet {
	l := d.Len()
	n := len(d.Fields)
	cols := make([]DataColumn, n)
	end := e + 1
	if end > l {
		end = l
	}
	for i := 0; i < n; i++ {
		cols[i] = d.Columns[i][b:end]
	}
	return DataSet{Fields: d.Fields, Types: d.Types, Columns: cols}
}

// MapRowAt datarow
func (d *DataSet) MapRowAt(i int) MapRow {
	n := len(d.Columns)
	if n == 0 || i >= len(d.Columns[0]) {
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
func (d *DataSet) AddDataRow(row Row) error {
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

// FloatAt index
func (d *DataSet) FloatAt(i int, field string) (float64, error) {
	k := d.FieldI(field)
	if k < 0 {
		return 0, ErrNoField
	}
	v := d.Columns[k][i]
	switch v.(type) {
	case float64:
		return v.(float64), nil
	case int64:
		return float64(v.(int64)), nil
	case bool:
		if v.(bool) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ErrType
	}
}

// Int64At index
func (d *DataSet) Int64At(i int, field string) (int64, error) {
	k := d.FieldI(field)
	if k < 0 {
		return 0, ErrNoField
	}
	v := d.Columns[k][i]
	switch v.(type) {
	case int64:
		return v.(int64), nil
	case float64:
		return int64(v.(float64)), nil
	case bool:
		if v.(bool) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ErrType
	}
}

// IntAt index
func (d *DataSet) IntAt(i int, field string) (int, error) {
	v, err := d.Int64At(i, field)
	return int(v), err
}

// BoolAt index
func (d *DataSet) BoolAt(i int, field string) (bool, error) {
	k := d.FieldI(field)
	if k < 0 {
		return false, ErrNoField
	}
	v := d.Columns[k][i]
	switch v.(type) {
	case []byte:
		val := v.([]byte)
		return len(val) > 0 && val[0] == 't', nil
	case int64:
		return v.(int64) == 1, nil
	case float64:
		return v.(float64) == 1, nil
	default:
		return false, ErrType
	}
}

// StrAt index
func (d *DataSet) StrAt(i int, field string) (string, error) {
	k := d.FieldI(field)
	if k < 0 {
		return "", ErrNoField
	}
	v := d.Columns[k][i]
	switch v.(type) {
	case []byte:
		return util.BytesToStr(v.([]byte)), nil
	default:
		return fmt.Sprint(v), nil
	}
}

// BytesAt index
func (d *DataSet) BytesAt(i int, field string) ([]byte, error) {
	k := d.FieldI(field)
	if k < 0 {
		return nil, ErrNoField
	}
	v := d.Columns[k][i]
	switch v.(type) {
	case []byte:
		return v.([]byte), nil
	default:
		return nil, ErrType
	}
}

// TimeAt index
func (d *DataSet) TimeAt(i int, field string) (time.Time, error) {
	k := d.FieldI(field)
	if k < 0 {
		return EmptyTime, ErrNoField
	}
	v := d.Columns[k][i]
	switch v.(type) {
	case time.Time:
		return v.(time.Time), nil
	default:
		return EmptyTime, ErrType
	}
}

// Int64sAt index
func (d *DataSet) Int64sAt(i int, field string) ([]int64, error) {
	k := d.FieldI(field)
	if k < 0 {
		return nil, ErrNoField
	}

	v := d.Columns[k][i]
	switch v.(type) {
	case []byte:
		return Current().Driver.Int64s(v.([]byte))

	default:
		return nil, ErrType
	}
}

// IntsAt index
func (d *DataSet) IntsAt(i int, field string) ([]int, error) {
	k := d.FieldI(field)
	if k < 0 {
		return nil, ErrNoField
	}

	v := d.Columns[k][i]
	switch v.(type) {
	case []byte:
		return Current().Driver.Ints(v.([]byte))

	default:
		return nil, ErrType
	}
}

// FloatsAt index
func (d *DataSet) FloatsAt(i int, field string) ([]float64, error) {
	k := d.FieldI(field)
	if k < 0 {
		return nil, ErrNoField
	}

	v := d.Columns[k][i]
	switch v.(type) {
	case []byte:
		return Current().Driver.Floats(v.([]byte))

	default:
		return nil, ErrType
	}
}

// StrsAt index
func (d *DataSet) StrsAt(i int, field string) ([]string, error) {
	k := d.FieldI(field)
	if k < 0 {
		return nil, ErrNoField
	}
	v := d.Columns[k][i]
	switch v.(type) {
	case []byte:
		return Current().Driver.Strings(v.([]byte))

	default:
		return nil, ErrType
	}
}

// ParseJSONAt index
func (d *DataSet) ParseJSONAt(i int, field string, vo interface{}) error {
	k := d.FieldI(field)
	if k < 0 {
		return ErrNoField
	}
	v := d.Columns[k][i]
	switch v.(type) {
	case []byte:
		return json.Unmarshal(v.([]byte), v)

	default:
		return ErrType
	}
}
