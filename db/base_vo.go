package db

import (
	"database/sql"
	"reflect"

	"github.com/kere/gno/libs/util"
)

// autotime:"true"
// skip: all / update /insert
// skipempty: all / update / insert

// IVO interface
type IVO interface {
	Table() string
}

// VO2InsertMapRow convert to MapRow
func VO2InsertMapRow(vo IVO) MapRow {
	cv := NewConvert(vo)
	return cv.ToMapRow(ActionInsert)
}

// VO2UpdateMapRow convert to MapRow
func VO2UpdateMapRow(vo IVO) MapRow {
	cv := NewConvert(vo)
	return cv.ToMapRow(ActionUpdate)
}

// VOCreate func
func VOCreate(vo IVO) error {
	ins := InsertBuilder{}
	row := VO2InsertMapRow(vo)
	_, err := ins.Table(vo.Table()).Insert(row)
	return err
}

// TxVOCreate func
func TxVOCreate(tx *Tx, vo IVO) error {
	ins := InsertBuilder{}
	row := VO2InsertMapRow(vo)
	_, err := ins.Table(vo.Table()).TxInsert(tx, row)
	tx.DoError(err)
	return err
}

// TxVOCreateAndReturnID create and return id
func TxVOCreateAndReturnID(tx *Tx, vo IVO) (sql.Result, error) {
	ins := InsertBuilder{}
	row := VO2InsertMapRow(vo)
	r, err := ins.Table(vo.Table()).ReturnID().TxInsert(tx, row)
	tx.DoError(err)
	return r, err
}

// VOCreateIfNotFound insert data if not found
// return true if insert
func VOCreateIfNotFound(vo IVO, where string, params ...interface{}) (bool, error) {
	e := ExistsBuilder{}
	if e.Table(vo.Table()).Where(where, params...).Exists() {
		return false, nil
	}

	return true, VOCreate(vo)
}

// TxVOCreateIfNotFound insert data if not found
// return true if insert
func TxVOCreateIfNotFound(tx *Tx, vo IVO, where string, params ...interface{}) (bool, error) {
	e := ExistsBuilder{}
	if exists, err := e.Table(vo.Table()).Where(where, params...).TxExists(tx); exists || err != nil {
		return false, err
	}

	return true, TxVOCreate(tx, vo)
}

// VOUpdate func
func VOUpdate(vo IVO, where string, params ...interface{}) error {
	u := UpdateBuilder{table: vo.Table()}
	row := VO2UpdateMapRow(vo)
	_, err := u.Where(where, params...).Update(row)
	return err
}

// TxVOUpdate update vo
func TxVOUpdate(tx *Tx, vo IVO, where string, params ...interface{}) error {
	u := UpdateBuilder{table: vo.Table()}
	row := VO2UpdateMapRow(vo)
	_, err := u.Where(where, params...).TxUpdate(tx, row)
	tx.DoError(err)
	return err
}

// VOUpdateByFields by list fields
func VOUpdateByFields(vo IVO, fields []string, where string, params ...interface{}) error {
	typ := reflect.TypeOf(vo)
	val := reflect.ValueOf(vo)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	l := typ.NumField()
	row := MapRow{}
	for i := 0; i < l; i++ {
		name := typ.Field(i).Tag.Get("json")
		if name != "" || util.InStrings(name, fields) {
			row[name] = val.Field(i).Interface()
		}
	}

	u := UpdateBuilder{table: vo.Table()}
	_, err := u.Where(where, params...).Update(row)
	return err
}

// VODelete func
func VODelete(vo IVO, where string, params ...interface{}) error {
	d := DeleteBuilder{table: vo.Table()}
	_, err := d.Where(where, params...).Delete()
	return err
}

// TxVODelete func
func TxVODelete(tx *Tx, vo IVO, where string, params ...interface{}) error {
	d := DeleteBuilder{table: vo.Table()}
	_, err := d.Where(where, params...).TxDelete(tx)
	tx.DoError(err)
	return err
}
