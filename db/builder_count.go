package db

import (
	"github.com/valyala/bytebufferpool"
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
	buf := bytebufferpool.Get()
	// driver := database.Driver
	buf.WriteString("SELECT count(1) as count FROM ")
	buf.WriteString(c.table)

	var rows MapRows
	var err error

	if c.isPrepare {
		buf.Write(bSQLWhere)
		buf.Write(c.GetDatabase().Driver.Adapt(cond, 0))
		rows, err = c.GetDatabase().QueryRowsPrepare(buf.String(), args...)
	} else {
		rows, err = c.GetDatabase().QueryRows(buf.String())
	}
	// bytePool.Put(buf)
	bytebufferpool.Put(buf)

	if err != nil {
		return -1, err
	}
	return rows[0].Int64(FieldCount), nil
}

//TxCount transaction Count db
func (c *CounterBuilder) TxCount(tx *Tx, cond string, args ...interface{}) (int64, error) {
	// s := bytes.Buffer{}
	buf := bytebufferpool.Get()

	// driver := database.Driver
	buf.WriteString("SELECT count(1) as count FROM ")
	buf.WriteString(c.table)

	var row MapRow
	var err error

	if c.isPrepare {
		buf.Write(bSQLWhere)
		buf.Write(c.GetDatabase().Driver.Adapt(cond, 0))
		row, err = tx.QueryOnePrepare(buf.String(), args...)
	} else {
		row, err = tx.QueryOne(buf.String())
	}

	if err != nil {
		return -1, err
	}

	return row.Int64(FieldCount), nil
}
