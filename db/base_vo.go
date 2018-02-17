package db

import (
	"encoding/json"
	"fmt"
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
	QueryOne(params ...interface{}) (DataRow, error)
	Query(params ...interface{}) (DataSet, error)
}

// BaseVO class
type BaseVO struct {
	target    IVO
	converter *StructConverter

	Fields string
	Table  string
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
	ins := NewInsertBuilder(b.Table)
	_, err := ins.Insert(b.target)

	if err != nil {
		return err
	}

	return nil
}

// CreateIfNotFound insert data if not found
// return true if insert
func (b *BaseVO) CreateIfNotFound(where string, params ...interface{}) (bool, error) {
	if NewExistsBuilder(b.Table).Where(where, params...).Exists() {
		return false, nil
	}

	return true, b.Create()
}

// Update func
func (b *BaseVO) Update(where string, params ...interface{}) error {
	_, err := NewUpdateBuilder(b.Table).Where(where, params...).Update(b.target)
	return err
}

// Delete func
func (b *BaseVO) Delete(where string, params ...interface{}) error {
	_, err := NewDeleteBuilder(b.Table).Where(where, params...).Delete()
	return err
}

// Query func
func (b *BaseVO) Query(params ...interface{}) (DataSet, error) {
	q := NewQueryBuilder(b.Table)
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
	q := NewQueryBuilder(b.Table)
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
