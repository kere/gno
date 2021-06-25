package dba

import (
	"database/sql"
)

type Builder struct {
	table string

	tx       *sql.Tx
	database *Database

	isTx      bool
	LastError error
}

// SetTable return string
func (b *Builder) SetTable(t string) {
	b.table = t
}

// GetTable return string
func (b *Builder) GetTable() string {
	return b.table
}

// GetTx tx
func (b *Builder) GetTx() *sql.Tx {
	return b.tx
}

// cQuery db tx
func (b *Builder) cQuery(isPool, isPrepare bool, sqlstr string, args []interface{}) (DataSet, error) {
	var err error
	var rows *sql.Rows
	// b.GetDatabase().Log(sqlstr, args)

	if isPrepare {
		var st *sql.Stmt
		if b.isTx {
			st, err = b.tx.Prepare(sqlstr)
		} else {
			st, err = b.database.DB().Prepare(sqlstr)
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
		if b.isTx {
			rows, err = b.tx.Query(sqlstr, args...)
		} else {
			rows, err = b.database.DB().Query(sqlstr, args...)
		}
		defer rows.Close()
		if err != nil {
			return EmptyDataSet, err
		}
	}

	return ScanToDataSet(rows, isPool)
}

// Exec db
func (b *Builder) Exec(sqlstr string, args []interface{}) (sql.Result, error) {
	// sqlstr := string(b.database.AdaptSql(item.Sql))
	// b.GetDatabase().Log(sqlstr, args)
	if b.isTx {
		return b.tx.Exec(sqlstr, args...)
	}

	return b.database.DB().Exec(sqlstr, args...)
}

// LastInsertID return lastid
func (b *Builder) LastInsertID(table, pkey string) int64 {
	var r *sql.Row
	if b.isTx {
		r = b.tx.QueryRow(Current().Driver.LastInsertID(table, pkey))
	} else {
		r = b.database.DB().QueryRow(Current().Driver.LastInsertID(table, pkey))
	}

	var count int64
	err := r.Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

// ExecPrepare db
func (b *Builder) ExecPrepare(sqlstr string, args []interface{}) (sql.Result, error) {
	// b.GetDatabase().Log(sqlstr, args)
	var st *sql.Stmt
	var err error
	if b.isTx {
		st, err = b.tx.Prepare(sqlstr)
	} else {
		st, err = b.database.DB().Prepare(sqlstr)
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
