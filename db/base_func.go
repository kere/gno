package db

import (
	"database/sql"
)

// Create f
func Create(table string, row MapRow) error {
	ins := InsertBuilder{}
	_, err := ins.Table(table).Insert(row)
	return err
}

// TxCreate func
func TxCreate(tx *Tx, table string, row MapRow) error {
	ins := InsertBuilder{}
	_, err := ins.Table(table).TxInsert(tx, row)

	return err
}

// TxCreateAndReturnID func
func TxCreateAndReturnID(tx *Tx, table string, row MapRow) (sql.Result, error) {
	ins := InsertBuilder{}
	return ins.Table(table).ReturnID().TxInsert(tx, row)
}

// CreateIfNotFound insert data if not found
// return true if insert
func CreateIfNotFound(table string, row MapRow, where string, params ...interface{}) (bool, error) {
	e := ExistsBuilder{}
	if e.Table(table).Where(where, params...).Exists() {
		return false, nil
	}

	return true, Create(table, row)
}

// Update func
func Update(table string, row MapRow, where string, params ...interface{}) error {
	u := UpdateBuilder{table: table}
	_, err := u.Where(where, params...).Update(row)
	return err
}

// TxUpdate func
func TxUpdate(tx *Tx, table string, row MapRow, where string, params ...interface{}) error {
	u := UpdateBuilder{table: table}
	_, err := u.Where(where, params...).TxUpdate(tx, row)
	return err
}

// Delete func
func Delete(table string, where string, params ...interface{}) error {
	d := DeleteBuilder{}
	_, err := d.Table(table).Where(where, params...).Delete()
	return err
}

// TxDelete func
func TxDelete(tx *Tx, table string, where string, params ...interface{}) error {
	d := DeleteBuilder{}
	_, err := d.Table(table).Where(where, params...).TxDelete(tx)
	return err
}
