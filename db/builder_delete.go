package db

import (
	"database/sql"

	"github.com/valyala/bytebufferpool"
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

func parseDelete(d *DeleteBuilder) string {
	// buf := bytePool.Get()
	buf := bytebufferpool.Get()
	// s := bytes.Buffer{}
	driver := Current().Driver
	buf.Write(bSQLDelete)
	buf.Write(bSQLFrom)
	buf.Write(driver.QuoteIdentifierB(d.table))
	if d.where != "" {
		buf.Write(bSQLWhere)
		buf.WriteString(d.GetDatabase().Driver.Adapt(d.where, 0))
	}
	str := buf.String()
	// bytePool.Put(buf)
	bytebufferpool.Put(buf)

	return str
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
	return d.GetDatabase().ExecPrepare(parseDelete(d), d.args...)
}

// TxDelete trunsaction delete
func (d *DeleteBuilder) TxDelete(tx *Tx) (sql.Result, error) {
	return tx.ExecPrepare(parseDelete(d), d.args...)
}
