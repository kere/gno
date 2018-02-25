package db

import (
	"bytes"
	"database/sql"
)

// Delete builder
type DeleteBuilder struct {
	builder
	table string
	where *CondParams
}

func NewDeleteBuilder(t string) *DeleteBuilder {
	return (&DeleteBuilder{}).Table(t)
}

func (d *DeleteBuilder) Table(t string) *DeleteBuilder {
	d.table = t
	return d
}

func (d *DeleteBuilder) parse() ([]byte, []interface{}) {
	var args []interface{}

	s := bytes.Buffer{}
	driver := Current().Driver
	s.Write(bSQLDelete)
	s.Write(bSQLFrom)
	s.WriteString(driver.QuoteField(d.table))
	if d.where != nil {
		s.Write(bSQLWhere)
		s.WriteString(d.where.Cond)
		args = d.where.Args
	}
	return s.Bytes(), args
}

func (d *DeleteBuilder) Where(cond string, args ...interface{}) *DeleteBuilder {
	if cond == "" {
		return d
	}
	d.where = &CondParams{cond, args}
	return d
}

func (d *DeleteBuilder) SqlState() *SqlState {
	sql, args := d.parse()
	return NewSqlState(sql, args...)
}

func (d *DeleteBuilder) Delete() (sql.Result, error) {
	return d.getDatabase().ExecPrepare(d.SqlState())
}

func (d *DeleteBuilder) TxDelete(tx *Tx) (sql.Result, error) {
	return tx.ExecPrepare(d.SqlState())
}
