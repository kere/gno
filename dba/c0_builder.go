package dba

import (
	"database/sql"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/myerr"
)

type Builder struct {
	table string
	tx    *sql.Tx

	isTx      bool
	LastError error

	database *Database
}

// Table return string
func (q *QueryBuilder) Table(t string) *QueryBuilder {
	q.table = t
	return q
}

// GetTable return string
func (q *QueryBuilder) GetTable() string {
	return q.table
}

func (t *Builder) GetDatabase() *Database {
	if t.database != nil {
		return t.database
	}
	t.database = Current()
	return t.database
}

func (t *Builder) SetDatabase(d *Database) {
	t.database = d
}

// GetTx tx
func (t *Builder) GetTx() *sql.Tx {
	return t.tx
}

// Begin tx
func (t *Builder) Begin() error {
	tx, err := t.GetDatabase().DB().Begin()
	if err != nil {
		t.GetDatabase().log.Alert(err).Stack()
		return err
	}
	t.tx = tx
	t.isTx = true
	return nil
}

// cQuery db tx
func (t *Builder) cQuery(isPool, isPrepare bool, sqlstr string, args []interface{}) (DataSet, error) {
	var err error
	var rows *sql.Rows
	t.GetDatabase().Log(sqlstr, args)

	if isPrepare {
		var st *sql.Stmt
		if t.isTx {
			st, err = t.tx.Prepare(sqlstr)
		} else {
			st, err = t.GetDatabase().DB().Prepare(sqlstr)
		}
		if err != nil {
			return EmptyDataSet, err
		}

		rows, err = st.Query(args...)
		defer rows.Close()
		if err != nil {
			return EmptyDataSet, err
		}
	} else {
		if t.isTx {
			rows, err = t.tx.Query(sqlstr, args...)
		} else {
			rows, err = t.GetDatabase().DB().Query(sqlstr, args...)
		}
		defer rows.Close()
		if err != nil {
			return EmptyDataSet, err
		}
	}

	return ScanToDataSet(rows, isPool)
}

// Exec db
func (t *Builder) Exec(sqlstr string, args []interface{}) (sql.Result, error) {
	// sqlstr := string(t.database.AdaptSql(item.Sql))
	t.GetDatabase().Log(sqlstr, args)
	if t.isTx {
		return t.tx.Exec(sqlstr, args...)
	}

	return t.GetDatabase().DB().Exec(sqlstr, args...)
}

// LastInsertID return lastid
func (t *Builder) LastInsertID(table, pkey string) int64 {
	var r *sql.Row
	if t.isTx {
		r = t.tx.QueryRow(Current().Driver.LastInsertID(table, pkey))
	} else {
		r = t.GetDatabase().DB().QueryRow(Current().Driver.LastInsertID(table, pkey))
	}

	var count int64
	err := r.Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

// ExecPrepare db
func (t *Builder) ExecPrepare(sqlstr string, args []interface{}) (sql.Result, error) {
	t.GetDatabase().Log(sqlstr, args)
	var st *sql.Stmt
	var err error
	if t.isTx {
		st, err = t.tx.Prepare(sqlstr)
	} else {
		st, err = t.GetDatabase().DB().Prepare(sqlstr)
	}
	if err != nil {
		return nil, err
	}
	r, err := st.Exec(args...)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// End func
func (t *Builder) End() error {
	t.isTx = false
	t.LastError = nil
	return t.Commit()
}

// Commit func
func (t *Builder) Commit() error {
	err := t.tx.Commit()
	if err != nil {
		log.App.Alert(err)
	}
	return err
}

// DoError tx
func (t *Builder) DoError(err error) bool {
	if err != nil {
		myerr.New(err).Log().Stack()
		err2 := t.tx.Rollback()
		if err2 != nil {
			t.GetDatabase().log.Alert(err2, "rollback faield")
			return true
		}
		t.LastError = err
		return true
	}
	return false
}

// Rollback err
func (t *Builder) Rollback() error {
	return t.tx.Rollback()
}
