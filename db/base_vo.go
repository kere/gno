package db

import (
	"fmt"
	"reflect"
)

// IVO interface
type IVO interface {
	Init(IVO) error
	Table() string
	SetTaget(IVO) error
	Create() error
	CreateIfNotFound(where string, params ...interface{}) (bool, error)
	Update(where string, params ...interface{}) error
	Delete(where string, params ...interface{}) error
}

// BaseVO class
type BaseVO struct {
	target           IVO
	val              reflect.Value
	converter        *StructConverter
	structPrimaryKey string
	tablePrimaryKey  string
}

// ToDataRow convert to DataRow
// ctype = insert | update |
func (b *BaseVO) ToDataRow(ctype string) DataRow {
	return b.converter.Struct2DataRow(ctype)
}

// Init func
func (b *BaseVO) Init(target IVO) error {
	if target == nil {
		return fmt.Errorf("value object is nil")
	}

	if err := b.SetTaget(target); err != nil {
		return err
	}
	return nil
}

// Table return string
func (b *BaseVO) Table() string {
	return b.target.Table()
}

// func (b *BaseVO) ExistsBuilder() *ExistsBuilder {
// 	return NewExistsBuilder(b.Table())
// }

// SetTaget func
func (b *BaseVO) SetTaget(t IVO) error {
	b.target = t
	b.converter = NewStructConvert(t)
	b.val = b.converter.val
	if !b.val.IsValid() {
		return fmt.Errorf("value object is invalid")
	}
	if b.val.Kind() == reflect.Ptr {
		b.val = b.val.Elem()
	}
	return nil
}

// Create func
func (b *BaseVO) Create() error {
	if b.target == nil {
		panic("vo.target is nil")
	}
	ins := NewInsertBuilder(b.Table())
	_, err := ins.Insert(b.target)

	if err != nil {
		return err
	}

	return nil
}

// CreateIfNotFound insert data if not found
// return true if insert
func (b *BaseVO) CreateIfNotFound(where string, params ...interface{}) (bool, error) {
	if NewExistsBuilder(b.Table()).Where(where, params...).Exists() {
		return false, nil
	}

	return true, b.Create()
}

// Update func
func (b *BaseVO) Update(where string, params ...interface{}) error {
	_, err := NewUpdateBuilder(b.Table()).Where(where, params...).Update(b.target)
	return err
}

// Delete func
func (b *BaseVO) Delete(where string, params ...interface{}) error {
	_, err := NewDeleteBuilder(b.Table()).Where(where, params...).Delete()
	return err
}
