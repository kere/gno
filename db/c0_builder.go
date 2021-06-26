package db

import (
	"database/sql"
)

type Builder struct {
	table string

	tx       *sql.Tx
	database *Database

	isTx      bool
	LastError error
	isPrepare bool
}

// NewBuilder return
func NewBuilder(t string) Builder {
	return Builder{table: t}
}

// SetTable return string
func (b *Builder) SetTable(t string) {
	b.table = t
}

// GetTable return string
func (b *Builder) GetTable() string {
	return b.table
}

// GetDatabase return string
func (b *Builder) GetDatabase() *Database {
	if b.database == nil {
		b.database = Current()
	}
	return b.database
}

// GetTx tx
func (b *Builder) GetTx() *sql.Tx {
	return b.tx
}

// cQuery db tx
func (b *Builder) cQuery(isPool bool, sqlstr string, args []interface{}) (DataSet, error) {
	var err error
	var rows *sql.Rows
	// b.GetDatabase().Log(sqlstr, args)

	if b.isPrepare {
		var st *sql.Stmt
		if b.isTx {
			st, err = b.tx.Prepare(sqlstr)
		} else {
			st, err = b.GetDatabase().DB().Prepare(sqlstr)
		}
		if err != nil {
			return EmptyDataSet, err
		}

		rows, err = st.Query(args...)
	} else {
		if b.isTx {
			rows, err = b.tx.Query(sqlstr, args...)
		} else {
			rows, err = b.GetDatabase().DB().Query(sqlstr, args...)
		}
	}
	if err != nil {
		return EmptyDataSet, err
	}
	defer rows.Close()

	return ScanToDataSet(rows, isPool)
}

// Exec db
func (b *Builder) Exec(sqlstr string, args []interface{}) (sql.Result, error) {
	// sqlstr := string(b.database.AdaptSql(item.Sql))
	// b.GetDatabase().Log(sqlstr, args)
	if b.isTx {
		return b.tx.Exec(sqlstr, args...)
	}

	return b.GetDatabase().DB().Exec(sqlstr, args...)
}

// LastInsertID return lastid
func (b *Builder) LastInsertID(table, pkey string) int64 {
	var r *sql.Row
	if b.isTx {
		r = b.tx.QueryRow(Current().Driver.LastInsertID(table, pkey))
	} else {
		r = b.GetDatabase().DB().QueryRow(Current().Driver.LastInsertID(table, pkey))
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
		st, err = b.GetDatabase().DB().Prepare(sqlstr)
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

// ScanToDataSet db
func ScanToDataSet(rows *sql.Rows, isPool bool) (DataSet, error) {
	cols, err := rows.Columns()
	if err != nil {
		return EmptyDataSet, err
	}

	colsNum := len(cols)

	typs, err := rows.ColumnTypes()
	if err != nil {
		return EmptyDataSet, err
	}

	fields := make([]string, colsNum)
	typItems := make([]ColType, colsNum)
	for i := 0; i < colsNum; i++ {
		typItems[i] = NewColType(typs[i])
		fields[i] = typs[i].Name()
	}

	var result DataSet
	if isPool {
		result = GetDataSet(fields)
	} else {
		result.Fields = fields
		result.Columns = make([][]interface{}, colsNum)
	}
	result.Types = typItems

	// var row, tem []interface{}
	row := GetRow(colsNum)
	tem := GetRow(colsNum)
	defer PutRow(row)
	defer PutRow(tem)
	for i := 0; i < colsNum; i++ {
		tem[i] = &row[i]
	}

	for rows.Next() {
		if err = rows.Scan(tem...); err != nil {
			return result, err
		}
		result.AddRow(row)
	}

	return result, rows.Err()
}
