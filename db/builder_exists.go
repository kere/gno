package db

import (
	"bytes"
)

// ExistsBuilder class
type ExistsBuilder struct {
	table string
	where string
	args  []interface{}
	builder
	isPrepare bool
}

// NewExistsBuilder new
func NewExistsBuilder(t string) *ExistsBuilder {
	return &ExistsBuilder{table: t}
}

// Table f
func (e *ExistsBuilder) Table(t string) *ExistsBuilder {
	e.table = t
	return e
}

// SetIsPrepare prepare sql
func (e *ExistsBuilder) SetIsPrepare(v bool) *ExistsBuilder {
	e.isPrepare = v
	return e
}

// GetIsPrepare f
func (e *ExistsBuilder) GetIsPrepare(v bool) bool {
	return e.isPrepare
}

// Where statement
func (e *ExistsBuilder) Where(s string, args ...interface{}) *ExistsBuilder {
	if s == "" {
		return e
	}
	e.where = s
	e.args = args
	return e
}

func (e *ExistsBuilder) parse() string {
	s := bytes.Buffer{}
	s.WriteString("SELECT 1 as field FROM ")
	s.WriteString(Current().Driver.QuoteField(e.table))

	if e.where != "" {
		s.WriteString(" WHERE ")
		s.WriteString(e.GetDatabase().Driver.Adapt(e.where, 0))
	}

	s.WriteString(" LIMIT 1")
	return s.String()
}

// Exists db
func (e *ExistsBuilder) Exists() bool {
	r, err := e.GetDatabase().QueryPrepare(e.parse(), e.args...)
	if err != nil {
		e.GetDatabase().log.Alert(err).Stack()
		panic(err)
	}
	return len(r) > 0
}

// TxExists trunsaction
func (e *ExistsBuilder) TxExists(tx *Tx) (bool, error) {
	return tx.Exists(e.parse(), e.args...)
}
