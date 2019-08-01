package db

import (
	"bytes"
)

// CounterBuilder class
type CounterBuilder struct {
	table     string
	isPrepare bool
	builder
}

// NewCounter new
func NewCounter(t string) *CounterBuilder {
	return &CounterBuilder{table: t}
}

// Table string
func (c *CounterBuilder) Table(t string) *CounterBuilder {
	c.table = t
	return c
}

// SetPrepare prepare sql
func (c *CounterBuilder) SetPrepare(v bool) *CounterBuilder {
	c.isPrepare = v
	return c
}

// GetPrepare get
func (c *CounterBuilder) GetPrepare() bool {
	return c.isPrepare
}

// Count db
func (c *CounterBuilder) Count(cond string, args ...interface{}) (int64, error) {
	// s := bytes.Buffer{}
	buf := bytePool.Get()
	// driver := database.Driver
	buf.WriteString("SELECT count(1) as count FROM ")
	buf.WriteString(c.table)

	var rows MapRows
	var err error

	if c.isPrepare {
		buf.Write(bSQLWhere)
		buf.WriteString(c.GetDatabase().Driver.Adapt(cond, 0))
		rows, err = c.GetDatabase().QueryRowsPrepare(buf.String(), args...)
	} else {
		rows, err = c.GetDatabase().QueryRows(buf.String())
	}

	if err != nil {
		return -1, err
	}
	return rows[0].Int64(FieldCount), nil
}

//TxCount transaction Count db
func (c *CounterBuilder) TxCount(tx *Tx, cond string, args ...interface{}) (int64, error) {
	s := bytes.Buffer{}
	// driver := database.Driver
	s.WriteString("SELECT count(1) as count FROM ")
	s.WriteString(c.table)

	var row MapRow
	var err error

	if c.isPrepare {
		s.Write(bSQLWhere)
		s.WriteString(c.GetDatabase().Driver.Adapt(cond, 0))
		row, err = tx.QueryOnePrepare(s.String(), args...)
	} else {
		row, err = tx.QueryOne(s.String())
	}

	if err != nil {
		return -1, err
	}

	return row.Int64(FieldCount), nil
}
