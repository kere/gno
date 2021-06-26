package db

import (
	"database/sql"

	"github.com/valyala/bytebufferpool"
)

// ExistsBuilder class
type ExistsBuilder struct {
	Builder
	where string
	args  []interface{}
}

// NewExists new
func NewExists(t string) ExistsBuilder {
	e := ExistsBuilder{}
	e.table = t
	return e
}

const (
	sExistsSQL = "SELECT 1 AS n FROM "
	// countField = "n"
)

// Table string
func (e *ExistsBuilder) Table(t string) *ExistsBuilder {
	e.table = t
	return e
}

// SetPrepare prepare sql
func (d *ExistsBuilder) Prepare(v bool) *ExistsBuilder {
	d.isPrepare = v
	return d
}

// Where sql
func (d *ExistsBuilder) Where(s string, args ...interface{}) *ExistsBuilder {
	d.where = s
	d.args = args
	return d
}

// Exists db
func (e *ExistsBuilder) NotExists() bool {
	return !e.Exists()
}

// Exists db
func (e *ExistsBuilder) Exists() bool {
	var err error
	var r sql.Result
	if e.isPrepare {
		r, err = e.ExecPrepare(parseExists(e), e.args)
	} else {
		r, err = e.Exec(parseExists(e), e.args)
	}
	if err != nil {
		e.database.log.Alert(err).Stack()
		return false
	}
	if r == nil {
		return false
	}

	n, _ := r.RowsAffected()
	return n > 0
}

func parseExists(e *ExistsBuilder) string {
	buf := bytebufferpool.Get()
	buf.WriteString(sExistsSQL)

	driver := e.GetDatabase().Driver
	driver.WriteQuoteIdentifier(buf, e.table)

	if e.where != "" {
		buf.Write(bSQLWhere)
		buf.WriteString(e.where)
	}

	buf.Write(bSQLLimitOne)
	str := buf.String()
	bytebufferpool.Put(buf)

	return str
}
