package db

import (
	"database/sql"
)

// QueryOne f
func QueryOne(table string, where string, params ...interface{}) (DataRow, error) {
	return NewQueryBuilder(table).Where(where, params...).QueryOne()
}

// Query f
func Query(table string, where string, params ...interface{}) (DataSet, error) {
	return NewQueryBuilder(table).Where(where, params...).Query()
}

// Create f
func Create(table string, row DataRow) error {
	ins := NewInsertBuilder(table)
	_, err := ins.Insert(row)
	return err
}

// TxCreate func
func TxCreate(tx *Tx, table string, row DataRow) error {
	ins := NewInsertBuilder(table)
	_, err := ins.TxInsert(tx, row)

	return err
}

// TxCreateAndReturnID func
func TxCreateAndReturnID(tx *Tx, table string, row DataRow) (sql.Result, error) {
	ins := NewInsertBuilder(table).ReturnID()
	return ins.TxInsert(tx, row)
}

// CreateIfNotFound insert data if not found
// return true if insert
func CreateIfNotFound(table string, row DataRow, where string, params ...interface{}) (bool, error) {
	if NewExistsBuilder(table).Where(where, params...).Exists() {
		return false, nil
	}

	return true, Create(table, row)
}

// Update func
func Update(table string, row DataRow, where string, params ...interface{}) error {
	_, err := NewUpdateBuilder(table).Where(where, params...).Update(row)
	return err
}

// TxUpdate func
func TxUpdate(tx *Tx, table string, row DataRow, where string, params ...interface{}) error {
	_, err := NewUpdateBuilder(table).Where(where, params...).TxUpdate(tx, row)
	return err
}

// Delete func
func Delete(table string, where string, params ...interface{}) error {
	_, err := NewDeleteBuilder(table).Where(where, params...).Delete()
	return err
}

// TxDelete func
func TxDelete(tx *Tx, table string, where string, params ...interface{}) error {
	_, err := NewDeleteBuilder(table).Where(where, params...).TxDelete(tx)
	return err
}