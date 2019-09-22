package db

// Create f
func Create(table string, row MapRow) error {
	ins := InsertBuilder{table: table}
	_, err := ins.Insert(row)
	return err
}

// TxCreate func
func TxCreate(tx *Tx, table string, row MapRow) error {
	ins := InsertBuilder{table: table}
	_, err := ins.TxInsert(tx, row)

	return err
}

// TxCreateAndReturnID func
func TxCreateAndReturnID(tx *Tx, table string, row MapRow) (int64, error) {
	ins := InsertBuilder{table: table}
	r, err := ins.ReturnID().TxInsert(tx, row)
	if err != nil {
		return -1, err
	}
	return r.LastInsertId()
}

// CreateIfNotFound insert data if not found
// return true if insert
func CreateIfNotFound(table string, row MapRow, where string, params ...interface{}) (bool, error) {
	e := ExistsBuilder{table: table}
	if e.Where(where, params...).Exists() {
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
	d := DeleteBuilder{table: table}
	_, err := d.Where(where, params...).Delete()
	return err
}

// TxDelete func
func TxDelete(tx *Tx, table string, where string, params ...interface{}) error {
	d := DeleteBuilder{table: table}
	_, err := d.Where(where, params...).TxDelete(tx)
	return err
}
