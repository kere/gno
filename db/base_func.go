package db

import (
	"database/sql"
)

// QueryOne f
func QueryOne(table string, where string, params ...interface{}) (DataRow, error) {
	q := QueryBuilder{}
	return q.Table(table).Where(where, params...).QueryOne()
}

// Query f
func Query(table string, where string, params ...interface{}) (DataSet, error) {
	q := QueryBuilder{}
	return q.Table(table).Where(where, params...).Query()
}

// Create f
func Create(table string, row DataRow) error {
	ins := InsertBuilder{}
	_, err := ins.Table(table).Insert(row)
	return err
}

// TxCreate func
func TxCreate(tx *Tx, table string, row DataRow) error {
	ins := InsertBuilder{}
	_, err := ins.Table(table).TxInsert(tx, row)

	return err
}

// TxCreateAndReturnID func
func TxCreateAndReturnID(tx *Tx, table string, row DataRow) (sql.Result, error) {
	ins := InsertBuilder{}
	return ins.Table(table).ReturnID().TxInsert(tx, row)
}

// CreateIfNotFound insert data if not found
// return true if insert
func CreateIfNotFound(table string, row DataRow, where string, params ...interface{}) (bool, error) {
	e := NewExistsBuilder(table)
	if e.Where(where, params...).Exists() {
		return false, nil
	}

	return true, Create(table, row)
}

// Update func
func Update(table string, row DataRow, where string, params ...interface{}) error {
	u := NewUpdateBuilder(table)
	_, err := u.Where(where, params...).Update(row)
	return err
}

// TxUpdate func
func TxUpdate(tx *Tx, table string, row DataRow, where string, params ...interface{}) error {
	u := NewUpdateBuilder(table)
	_, err := u.Where(where, params...).TxUpdate(tx, row)
	return err
}

// Delete func
func Delete(table string, where string, params ...interface{}) error {
	d := NewDeleteBuilder(table)
	_, err := d.Where(where, params...).Delete()
	return err
}

// TxDelete func
func TxDelete(tx *Tx, table string, where string, params ...interface{}) error {
	d := NewDeleteBuilder(table)
	_, err := d.Where(where, params...).TxDelete(tx)
	return err
}
