package dba

import (
	"database/sql"

	"github.com/valyala/bytebufferpool"
)

// ExistsBuilder class
type ExistsBuilder struct {
	Builder
	IsPrepare bool
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

func parseExists(e *ExistsBuilder, where string) string {
	buf := bytebufferpool.Get()
	buf.WriteString(sExistsSQL)

	driver := e.GetDatabase().Driver
	driver.WriteQuoteIdentifier(buf, e.table)

	if where != "" {
		buf.Write(bSQLWhere)
		buf.WriteString(where)
	}

	buf.Write(bSQLLimitOne)
	str := buf.String()
	bytebufferpool.Put(buf)

	return str
}

// Exists db
func (e *ExistsBuilder) NotExists(where string, args ...interface{}) bool {
	return !e.Exists(where, args...)
}

// Exists db
func (e *ExistsBuilder) Exists(where string, args ...interface{}) bool {
	var err error
	var r sql.Result
	if e.IsPrepare {
		r, err = e.ExecPrepare(parseExists(e, where), args)
	} else {
		r, err = e.Exec(parseExists(e, where), args)
	}
	if err != nil {
		e.GetDatabase().log.Alert(err).Stack()
		return false
	}
	if r == nil {
		return false
	}

	n, _ := r.RowsAffected()
	return n > 0
}
