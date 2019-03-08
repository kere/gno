package db

import (
	"bytes"
	"database/sql"
)

// DeleteBuilder class
type DeleteBuilder struct {
	builder
	table string
	where string
	args  []interface{}
}

// NewDelete func
func NewDelete(t string) *DeleteBuilder {
	return &DeleteBuilder{table: t}
}

//Table string
func (d *DeleteBuilder) Table(t string) *DeleteBuilder {
	d.table = t
	return d
}

func (d *DeleteBuilder) parse() string {
	s := bytes.Buffer{}
	driver := Current().Driver
	s.Write(bSQLDelete)
	s.Write(bSQLFrom)
	s.WriteString(driver.QuoteField(d.table))
	if d.where != "" {
		s.Write(bSQLWhere)
		s.WriteString(d.GetDatabase().Driver.Adapt(d.where, 0))
	}
	return s.String()
}

// Where conditions
func (d *DeleteBuilder) Where(cond string, args ...interface{}) *DeleteBuilder {
	if cond == "" {
		return d
	}
	d.where = cond
	d.args = args
	return d
}

// Delete delete
func (d *DeleteBuilder) Delete() (sql.Result, error) {
	return d.GetDatabase().ExecPrepare(d.parse(), d.args...)
}

// TxDelete trunsaction delete
func (d *DeleteBuilder) TxDelete(tx *Tx) (sql.Result, error) {
	return tx.ExecPrepare(d.parse(), d.args...)
}
