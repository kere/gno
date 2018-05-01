package db

import (
	"bytes"
)

// Counter builder
type CounterBuilder struct {
	table string
	builder
}

func NewCounterBuilder(t string) *CounterBuilder {
	return (&CounterBuilder{}).Table(t)
}

func (c *CounterBuilder) Table(t string) *CounterBuilder {
	c.table = t
	return c
}

func (c *CounterBuilder) Count(cond string, args ...interface{}) (int64, error) {
	s := bytes.Buffer{}
	database := c.getDatabase()
	// driver := database.Driver
	s.WriteString("SELECT count(1) as count FROM ")
	s.WriteString(c.table)

	var r DataSet
	var err error

	if cond != "" {
		s.Write(bSQLWhere)
		s.WriteString(cond)
		r, err = database.QueryPrepare(NewSqlState(s.Bytes(), args...))
	} else {
		r, err = database.Query(NewSqlState(s.Bytes()))
	}

	if err != nil {
		return -1, err
	}
	return r[0].Int64("count"), nil
}

func (c *CounterBuilder) TxCount(tx *Tx, cond string, args ...interface{}) (int64, error) {
	s := bytes.Buffer{}
	// driver := database.Driver
	s.WriteString("SELECT count(1) as count FROM ")
	s.WriteString(c.table)

	var row DataRow
	var err error

	if cond != "" {
		s.Write(bSQLWhere)
		s.WriteString(cond)
		row, err = tx.QueryOne(NewSqlState(s.Bytes(), args...))
	} else {
		row, err = tx.QueryOne(NewSqlState(s.Bytes()))
	}

	if err != nil {
		return -1, err
	}

	return row.Int64("count"), nil
}
