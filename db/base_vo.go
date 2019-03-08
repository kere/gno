package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/kere/gno/libs/util"
)

// autotime:"true"
// skip: all / update /insert
// skipempty: all / update / insert

// IVO interface
type IVO interface {
	Init(string, IVO)
	GetTable() string
	ToDataRow(ctype string) DataRow
	ToJSON(ctype string) string
	Create() error
	CreateIfNotFound(where string, params ...interface{}) (bool, error)
	Update(where string, params ...interface{}) error
	Delete(where string, params ...interface{}) error
	Order(string) *QueryBuilder
	QueryOne(params ...interface{}) (DataRow, error)
	Query(params ...interface{}) (DataSet, error)
}

// BaseVO class
type BaseVO struct {
	target    IVO
	converter *StructConverter

	Fields     string
	Table      string
	querybuild *QueryBuilder
}

// ToDataRow convert to DataRow
// ctype = insert | update |
func (b *BaseVO) ToDataRow(ctype string) DataRow {
	if b.converter == nil {
		b.converter = NewStructConvert(b.target)
	}
	return b.converter.Struct2DataRow(ctype)
}

// ToJSON convert to DataRow
// ctype = insert | update |
func (b *BaseVO) ToJSON(ctype string) string {
	row := b.ToDataRow(ctype)
	src, _ := json.Marshal(row)
	return string(src)
}

// Init func
func (b *BaseVO) Init(table string, target IVO) {
	if target == nil {
		panic(fmt.Errorf("value object is nil"))
	}

	b.target = target
	b.Table = table
}

// GetTable return string
func (b *BaseVO) GetTable() string {
	return b.Table
}

// Create func
func (b *BaseVO) Create() error {
	if b.target == nil {
		panic("vo.target is nil")
	}
	ins := InsertBuilder{}
	_, err := ins.Table(b.Table).Insert(b.target)

	if err != nil {
		return err
	}

	return nil
}

// TxCreate func
func (b *BaseVO) TxCreate(tx *Tx) error {
	if b.target == nil {
		panic("vo.target is nil")
	}
	ins := InsertBuilder{}
	_, err := ins.Table(b.Table).TxInsert(tx, b.target)
	tx.DoError(err)
	return err
}

// TxCreateAndReturnID func
func (b *BaseVO) TxCreateAndReturnID(tx *Tx) (sql.Result, error) {
	if b.target == nil {
		panic("vo.target is nil")
	}
	ins := InsertBuilder{}
	r, err := ins.Table(b.Table).ReturnID().TxInsert(tx, b.target)
	tx.DoError(err)
	return r, err
}

// TxExists is found
func (b *BaseVO) TxExists(tx *Tx, where string, params ...interface{}) (bool, error) {
	e := ExistsBuilder{}
	isok, err := e.Table(b.Table).Where(where, params...).TxExists(tx)
	tx.DoError(err)
	return isok, err
}

// Exists is found
func (b *BaseVO) Exists(where string, params ...interface{}) bool {
	e := ExistsBuilder{}
	return e.Table(b.Table).Where(where, params...).Exists()
}

// CreateIfNotFound insert data if not found
// return true if insert
func (b *BaseVO) CreateIfNotFound(where string, params ...interface{}) (bool, error) {
	e := ExistsBuilder{}
	if e.Table(b.Table).Where(where, params...).Exists() {
		return false, nil
	}

	return true, b.Create()
}

// TxCreateIfNotFound insert data if not found
// return true if insert
func (b *BaseVO) TxCreateIfNotFound(tx *Tx, where string, params ...interface{}) (bool, error) {
	e := ExistsBuilder{}
	if exists, err := e.Table(b.Table).Where(where, params...).TxExists(tx); exists || tx.DoError(err) {
		return false, err
	}

	return true, b.TxCreate(tx)
}

// Update func
func (b *BaseVO) Update(where string, params ...interface{}) error {
	u := UpdateBuilder{table: b.Table}
	_, err := u.Where(where, params...).Update(b.target)
	return err
}

// TxUpdate func
func (b *BaseVO) TxUpdate(tx *Tx, where string, params ...interface{}) error {
	u := UpdateBuilder{table: b.Table}
	_, err := u.Where(where, params...).TxUpdate(tx, b.target)
	tx.DoError(err)
	return err
}

// UpdateFields func
func (b *BaseVO) UpdateFields(fields []string, where string, params ...interface{}) error {
	typ := reflect.TypeOf(b.target)
	val := reflect.ValueOf(b.target)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	l := typ.NumField()
	dat := DataRow{}
	for i := 0; i < l; i++ {
		name := typ.Field(i).Tag.Get("json")
		if name != "" || util.InStrings(name, fields) {
			dat[name] = val.Field(i).Interface()
		}
	}

	u := UpdateBuilder{table: b.Table}
	_, err := u.Where(where, params...).Update(dat)
	return err
}

// Delete func
func (b *BaseVO) Delete(where string, params ...interface{}) error {
	d := DeleteBuilder{table: b.Table}
	_, err := d.Where(where, params...).Delete()
	return err
}

// TxDelete func
func (b *BaseVO) TxDelete(tx *Tx, where string, params ...interface{}) error {
	d := DeleteBuilder{table: b.Table}
	_, err := d.Where(where, params...).TxDelete(tx)
	tx.DoError(err)
	return err
}

// getQueryBuilder func
func (b *BaseVO) getQueryBuilder() *QueryBuilder {
	if b.querybuild != nil {
		return b.querybuild
	}
	e := QueryBuilder{}
	b.querybuild = e.Table(b.Table)
	return b.querybuild
}

// Order func
func (b *BaseVO) Order(s string) *QueryBuilder {
	b.getQueryBuilder().Order(s)
	return b.querybuild
}

// Query func
func (b *BaseVO) Query(params ...interface{}) (DataSet, error) {
	q := b.getQueryBuilder()
	if len(params) == 1 {
		q.Where(fmt.Sprint(params[0]))
	} else if len(params) > 1 {
		q.Where(fmt.Sprint(params[0]), params[1:]...)
	}

	if b.Fields != "" {
		q.Select(b.Fields)
	}
	return q.Query()
}

// QueryOne func
func (b *BaseVO) QueryOne(params ...interface{}) (DataRow, error) {
	q := b.getQueryBuilder()
	if len(params) == 1 {
		q.Where(fmt.Sprint(params[0]))
	} else if len(params) > 1 {
		q.Where(fmt.Sprint(params[0]), params[1:]...)
	}

	if b.Fields != "" {
		q.Select(b.Fields)
	}

	return q.QueryOne()
}
