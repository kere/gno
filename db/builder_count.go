package db

import "bytes"

// Counter builder
type CounterBuilder struct {
	table string
	builder
}

func NewCounterBuilder(t string) *CounterBuilder {
	return (&CounterBuilder{}).Table(t)
}

func (this *CounterBuilder) Table(t string) *CounterBuilder {
	this.table = t
	return this
}

func (this *CounterBuilder) Count(cond string, args ...interface{}) (int64, error) {
	s := bytes.Buffer{}
	database := this.getDatabase()
	// driver := database.Driver
	s.WriteString("SELECT count(1) as count FROM ")
	s.WriteString(this.table)

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
