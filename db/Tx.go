package db

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/myerr"
)

//Tx class
type Tx struct {
	builder
	conn      *sql.DB
	tx        *sql.Tx
	IsError   bool
	LastError error
}

// NewTx tx
func NewTx() *Tx {
	return &Tx{}
}

// GetTx tx
func (t *Tx) GetTx() *sql.Tx {
	return t.tx
}

// Begin tx
func (t *Tx) Begin() error {
	t.conn = t.GetDatabase().Conn()
	tx, err := t.conn.Begin()
	if err != nil {
		t.IsError = true
		log.App.Alert(err).Stack()
		return err
	}
	t.tx = tx
	t.IsError = false
	return nil
}

// QueryOne tx
func (t *Tx) QueryOne(sql string, args ...interface{}) (MapRow, error) {
	r, err := t.QueryRows(sql, args...)
	if err != nil {
		return nil, err
	}

	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

// QueryOnePrepare tx
func (t *Tx) QueryOnePrepare(sql string, args ...interface{}) (MapRow, error) {
	r, err := t.QueryRowsPrepare(sql, args...)
	if err != nil {
		return nil, err
	}

	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

// Exists db
func (t *Tx) Exists(sql string, args ...interface{}) (bool, error) {
	r, err := t.QueryRows(sql, args...)
	if err != nil {
		return false, err
	}
	return len(r) > 0, nil
}

// Query db tx
func (t *Tx) Query(sqlstr string, args ...interface{}) (DataSet, error) {
	dataset, _, err := t.cQuery(0, 0, sqlstr, args...)
	return dataset, err
}

// QueryPrepare db tx
func (t *Tx) QueryPrepare(sqlstr string, args ...interface{}) (DataSet, error) {
	dataset, _, err := t.cQuery(0, 1, sqlstr, args...)
	return dataset, err
}

// QueryRows db tx
func (t *Tx) QueryRows(sqlstr string, args ...interface{}) (MapRows, error) {
	_, datarows, err := t.cQuery(1, 0, sqlstr, args...)
	return datarows, err
}

// QueryRowsPrepare db tx
func (t *Tx) QueryRowsPrepare(sqlstr string, args ...interface{}) (MapRows, error) {
	_, datarows, err := t.cQuery(1, 1, sqlstr, args...)
	return datarows, err
}

// cQuery db tx
func (t *Tx) cQuery(mode, qmode int, sqlstr string, args ...interface{}) (DataSet, MapRows, error) {
	var dataset DataSet
	if t.IsError {
		return dataset, nil, nil
	}
	var err error
	var rows *sql.Rows

	t.GetDatabase().Log(sqlstr, args)
	if mode == 1 {
		st, err := t.tx.Prepare(sqlstr)
		if err != nil {
			return dataset, nil, err
		}

		rows, err = st.Query(args...)
		defer rows.Close()
		if err != nil {
			return dataset, nil, err
		}

	} else {
		rows, err = t.tx.Query(sqlstr, args...)
		defer rows.Close()
		if err != nil {
			return dataset, nil, err
		}

	}

	var datarows MapRows
	if mode == 1 {
		datarows, err = ScanToMapRows(rows)
	} else {
		dataset, err = ScanToDataSet(rows)
	}

	return dataset, datarows, err
}

// Exec db
func (t *Tx) Exec(sqlstr string, args ...interface{}) (sql.Result, error) {
	if t.IsError {
		return nil, nil
	}
	// sqlstr := string(t.database.AdaptSql(item.Sql))
	t.GetDatabase().Log(sqlstr, args)

	r, err := t.tx.Exec(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// LastInsertID return lastid
func (t *Tx) LastInsertID(table, pkey string) int64 {
	if t.IsError {
		return -1
	}
	r := t.tx.QueryRow(Current().Driver.LastInsertID(table, pkey))

	var count int64
	err := r.Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

// ExecPrepare db
func (t *Tx) ExecPrepare(sqlstr string, args ...interface{}) (sql.Result, error) {
	if t.IsError {
		return nil, nil
	}
	t.GetDatabase().Log(sqlstr, args)
	// sqlstr := string(t.database.AdaptSql(item.Sql))
	st, err := t.tx.Prepare(sqlstr)
	if err != nil {
		return nil, err
	}
	r, err := st.Exec(args...)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// PGCopyIn db
func (t *Tx) PGCopyIn(table string, fields []string, rows []MapRow) (int, error) {
	l := len(rows)
	step := 200
	n := l/step + 1
	count := 0
	for i := 0; i < n; i++ {
		b := i * step
		if b >= l {
			break
		}
		e := b + step
		if e > l {
			e = l
		}
		err := t.pgCopyIn(table, fields, rows[b:e])
		if err != nil {
			return count, err
		}
		count += e - b
	}
	return count, nil
}

func (t *Tx) pgCopyIn(table string, fields []string, rows []MapRow) error {
	stmt, err := t.tx.Prepare(pq.CopyIn(table, fields...))
	if err != nil {
		return err
	}
	n := len(fields)
	l := len(rows)
	for i := 0; i < l; i++ {
		arr := make([]interface{}, n)
		for k := 0; k < n; k++ {
			arr[k] = rows[i][fields[k]]
		}
		_, err := stmt.Exec(arr...)
		if err != nil {
			return err
		}
	}
	_, err = stmt.Exec()
	return err
}

// func (t *Tx) close() error {
// 	if t.conn == nil {
// 		return nil
// 	}
// 	return t.conn.Close()
// }

// End func
func (t *Tx) End() error {
	return t.Commit()
	// return t.close()
}

// Commit func
func (t *Tx) Commit() error {
	if t.IsError {
		return nil
	}
	err := t.tx.Commit()
	if err != nil {
		log.App.Alert(err)
	}
	return err
}

// DoError tx
func (t *Tx) DoError(err error) bool {
	if t.IsError {
		return true
	}

	if err != nil {
		myerr.New(err).Log().Stack()
		err2 := t.tx.Rollback()
		if err2 != nil {
			log.App.Alert(err2, "rollback faield")
			return true
		}
		t.IsError = true
		t.LastError = err
		return true
	}
	return false
}

// Rollback err
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
