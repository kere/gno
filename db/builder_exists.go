package db

import (
	"bytes"
)

// Query builder
type ExistsBuilder struct {
	table string
	where *CondParams
	builder
	isPrepare bool
}

func NewExistsBuilder(t string) *ExistsBuilder {
	return (&ExistsBuilder{}).Table(t)
}

func (e *ExistsBuilder) Table(t string) *ExistsBuilder {
	e.table = t
	return e
}

func (e *ExistsBuilder) IsPrepare(v bool) *ExistsBuilder {
	e.isPrepare = v
	return e
}

func (e *ExistsBuilder) GetIsPrepare(v bool) bool {
	return e.isPrepare
}

func (e *ExistsBuilder) Where(s string, args ...interface{}) *ExistsBuilder {
	if s == "" {
		return e
	}
	e.where = &CondParams{s, args}
	return e
}

func (e *ExistsBuilder) SqlState() *SqlState {
	return NewSqlState(e.parse())
}

func (e *ExistsBuilder) parse() ([]byte, []interface{}) {
	var args []interface{}
	s := bytes.Buffer{}
	s.WriteString("SELECT 1 as field FROM ")
	s.WriteString(Current().Driver.QuoteField(e.table))

	if e.where != nil {
		s.WriteString(" WHERE ")
		s.WriteString(e.where.Cond)
		args = e.where.Args
	}

	s.WriteString(" LIMIT 1")
	return s.Bytes(), args
}

func (e *ExistsBuilder) Exists() bool {
	r, err := e.getDatabase().QueryPrepare(e.SqlState())

	if err != nil {
		panic(err)
	} else {
		return len(r) > 0
	}
}

func (e *ExistsBuilder) TxExists(tx *Tx) (bool, error) {
	return tx.Exists(e.SqlState())
}
