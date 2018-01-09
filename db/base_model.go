package db

import (
	"database/sql"
	"fmt"
)

type IModel interface {
	Init(IVO) error
	QueryBuilder() *QueryBuilder
	SetTablePrimarykeyName(string)
	TablePrimarykeyName() string
	Table() string
	EnableCache(bool)

	QueryByID(int64) (DataRow, error)
	FindByID(int64) (IVO, error)

	QueryOne(string, ...interface{}) (DataRow, error)
	FindOne(string, ...interface{}) (IVO, error)

	Query(string, ...interface{}) (DataSet, error)
	Find(string, ...interface{}) (VODataSet, error)

	Delete(int64) (sql.Result, error)
	ClearCache(int64)
}

type BaseModel struct {
	vo       IVO
	table    string
	pkeyName string
	query    *QueryBuilder
}

// initalize database model
// vo : model value object, it must be base on BaseVO.
func (b *BaseModel) Init(vo IVO) error {
	var err error
	err = vo.Init(vo)
	b.vo = vo
	if err != nil {
		return err
	}
	b.pkeyName = "id"

	b.EnableCache(false)
	return nil
}

// get table primary key name, default primary key name is "id".
func (b *BaseModel) QueryBuilder() *QueryBuilder {
	if b.query == nil {
		b.query = NewQueryBuilder(b.Table()).Struct(b.vo)
	}
	return b.query
}

// set table primary key name, default pkey name is "id".
func (b *BaseModel) SetTablePrimarykeyName(t string) {
	b.pkeyName = t
}

// get table primary key name, default primary key name is "id".
func (b *BaseModel) TablePrimarykeyName() string {
	return b.pkeyName
}

// get table name
func (b *BaseModel) Table() string {
	return b.vo.Table()
}

// enable cache
func (b *BaseModel) EnableCache(c bool) {
	if c {
		b.QueryBuilder().Cache()
	} else {
		b.QueryBuilder().DisableCache()
	}
}

// query db record by primary key limit one, and return DataRow
func (b *BaseModel) QueryByID(id int64) (DataRow, error) {
	if id < 1 {
		return nil, fmt.Errorf("primary key mest be > 0")
	}
	return b.QueryBuilder().Where(fmt.Sprint(b.pkeyName, "=?"), id).QueryOne()
}

// find db record by primary key limit one, and return interface{} of value object
func (b *BaseModel) FindByID(id int64) (IVO, error) {
	if b.vo == nil {
		return nil, fmt.Errorf("model vaule object not be set.")
	}
	return b.QueryBuilder().Where(fmt.Sprint(b.pkeyName, "=?"), id).FindOne()
}

// query db record limit one, and return DataRow
func (b *BaseModel) QueryOne(cond string, args ...interface{}) (DataRow, error) {
	return b.QueryBuilder().Where(cond, args...).QueryOne()
}

// find db record limit one, and return interface{} of value object
func (b *BaseModel) FindOne(cond string, args ...interface{}) (IVO, error) {
	return b.QueryBuilder().Where(cond, args...).FindOne()
}

// query db record, and return DataSet
func (b *BaseModel) Query(cond string, args ...interface{}) (DataSet, error) {
	return b.QueryBuilder().Where(cond, args...).Query()
}

// find db record, and return VODataSet
func (b *BaseModel) Find(cond string, args ...interface{}) (VODataSet, error) {
	if b.vo == nil {
		return nil, fmt.Errorf("model vaule object not be set.")
	}
	return b.QueryBuilder().Where(cond, args...).Find()
}

// query db record pagination, and return DataSet
func (b *BaseModel) PageQuery(page, pageSize int, cond string, args ...interface{}) (DataSet, error) {
	return b.QueryBuilder().Page(page, pageSize).Where(cond, args...).Query()
}

// find db record, and return VODataSet
func (b *BaseModel) PageFind(page, pageSize int, cond string, args ...interface{}) (VODataSet, error) {
	if b.vo == nil {
		return nil, fmt.Errorf("model vaule object not be set.")
	}
	return b.QueryBuilder().Page(page, pageSize).Where(cond, args...).Find()
}

// delete by primary key
func (b *BaseModel) Delete(id int64) (sql.Result, error) {
	if id < 1 {
		return nil, fmt.Errorf("primary key mest be > 0")
	}
	return NewDeleteBuilder(b.table).Where(fmt.Sprint(b.pkeyName, "=?"), id).Delete()
}

// delete by primary key
func (b *BaseModel) ClearCache(id int64) {
	b.QueryBuilder().Where(fmt.Sprint(b.pkeyName, "=?"), id).ClearCache()
}
