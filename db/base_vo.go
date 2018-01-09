package db

import (
	"fmt"
	"reflect"
)

// IVO interface
type IVO interface {
	Init(IVO) error
	SetStructPrimaryKeyName(string)
	Table() string
	SetTaget(IVO) error
	Create() error
	CreateIfNotFound(where string, params ...interface{}) (bool, error)
	Update(int64) error
	UpdateWhere(where string, params ...interface{}) error
	Delete() error
	DeleteWhere(where string, params ...interface{}) error
}

// BaseVO class
type BaseVO struct {
	target           IVO
	val              reflect.Value
	converter        *StructConverter
	table            string
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

	b.SetStructPrimaryKeyName("ID")

	return nil
}

// Table return string
func (b *BaseVO) Table() string {
	if b.target != nil {
		return b.target.Table()
	}

	return ""
}

// func (b *BaseVO) ExistsBuilder() *ExistsBuilder {
// 	return NewExistsBuilder(b.Table())
// }

// SetStructPrimaryKeyName func
func (b *BaseVO) SetStructPrimaryKeyName(t string) {
	b.structPrimaryKey = t
	if structfield, found := b.val.Type().FieldByName(b.structPrimaryKey); found {
		b.tablePrimaryKey = structfield.Tag.Get("json")
	}

	if b.tablePrimaryKey == "" {
		b.tablePrimaryKey = "id"
	}
}

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
func (b *BaseVO) Update(id int64) error {
	if b.tablePrimaryKey == "" {
		return fmt.Errorf("use init() to initialize VO")
	}

	if id < 1 {
		return fmt.Errorf("update failed, id must be >0")
	}

	_, err := NewUpdateBuilder(b.Table()).Where(b.tablePrimaryKey+"=?", id).Update(b.target)
	return err
}

// UpdateWhere func
func (b *BaseVO) UpdateWhere(where string, params ...interface{}) error {
	_, err := NewUpdateBuilder(b.Table()).Where(where, params...).Update(b.target)
	return err
}

// Delete func
func (b *BaseVO) Delete() error {
	id := b.val.FieldByName("Id").Int()
	if id < 1 {
		return fmt.Errorf("update failed, id must be >0")
	}

	_, err := NewDeleteBuilder(b.Table()).Where("id=?", id).Delete()
	return err
}

// DeleteWhere func
func (b *BaseVO) DeleteWhere(where string, params ...interface{}) error {
	_, err := NewDeleteBuilder(b.Table()).Where(where, params...).Delete()
	return err
}
