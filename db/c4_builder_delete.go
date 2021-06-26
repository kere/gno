package db

import (
	"database/sql"

	"github.com/valyala/bytebufferpool"
)

// DeleteBuilder class
type DeleteBuilder struct {
	Builder
	where string
	args  []interface{}
}

// NewDelete func
func NewDelete(t string) DeleteBuilder {
	d := DeleteBuilder{}
	d.table = t
	return d
}

// Table string
func (d *DeleteBuilder) Table(t string) *DeleteBuilder {
	d.table = t
	return d
}

// SetPrepare prepare sql
func (d *DeleteBuilder) Prepare(v bool) *DeleteBuilder {
	d.isPrepare = v
	return d
}

// Where sql
func (d *DeleteBuilder) Where(s string, args ...interface{}) *DeleteBuilder {
	d.where = s
	d.args = args
	return d
}

// Delete delete
func (d *DeleteBuilder) Delete() (sql.Result, error) {
	if d.isPrepare {
		return d.ExecPrepare(parseDelete(d), d.args)
	}
	return d.Exec(parseDelete(d), d.args)
}

func parseDelete(d *DeleteBuilder) string {
	buf := bytebufferpool.Get()
	buf.Write(bSQLDelete)
	buf.Write(bSQLFrom)
	driver := d.GetDatabase().Driver
	driver.WriteQuoteIdentifier(buf, d.table)
	if d.where != "" {
		buf.Write(bSQLWhere)
		buf.WriteString(d.where)
	}
	str := buf.String()
	bytebufferpool.Put(buf)
	return str
}
