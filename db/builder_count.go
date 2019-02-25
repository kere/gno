package db

import (
	"bytes"
)

// CounterBuilder class
type CounterBuilder struct {
	table string
	builder
}

// NewCounterBuilder new
func NewCounterBuilder(t string) CounterBuilder {
	return CounterBuilder{table: t}
}

// Table string
func (c *CounterBuilder) Table(t string) *CounterBuilder {
	c.table = t
	return c
}

// Count db
func (c *CounterBuilder) Count(cond string, args ...interface{}) (int64, error) {
	s := bytes.Buffer{}
	// driver := database.Driver
	s.WriteString("SELECT count(1) as count FROM ")
	s.WriteString(c.table)

	var r DataSet
	var err error

	if cond != "" {
		s.Write(bSQLWhere)
		s.WriteString(c.GetDatabase().Driver.Adapt(cond, 0))
		r, err = c.GetDatabase().QueryPrepare(s.String(), args...)
	} else {
		r, err = c.GetDatabase().Query(s.String())
	}

	if err != nil {
		return -1, err
	}
	return r[0].Int64(FieldCount), nil
}

//TxCount transaction Count db
func (c *CounterBuilder) TxCount(tx *Tx, cond string, args ...interface{}) (int64, error) {
	s := bytes.Buffer{}
	// driver := database.Driver
	s.WriteString("SELECT count(1) as count FROM ")
	s.WriteString(c.table)

	var row DataRow
	var err error

	if cond != "" {
		s.Write(bSQLWhere)
		s.WriteString(c.GetDatabase().Driver.Adapt(cond, 0))
		row, err = tx.QueryOne(s.String(), args...)
	} else {
		row, err = tx.QueryOne(s.String())
	}

	if err != nil {
		return -1, err
	}

	return row.Int64(FieldCount), nil
}
