package db

import (
	"bytes"
	"database/sql"
)

// Update builder
type UpdateBuilder struct {
	table string
	where *CondParams
	builder
}

func NewUpdateBuilder(t string) *UpdateBuilder {
	return (&UpdateBuilder{}).Table(t)
}

func (u *UpdateBuilder) Table(t string) *UpdateBuilder {
	u.table = t
	return u
}

func (u *UpdateBuilder) parse(data interface{}) ([]byte, []interface{}) {
	keys, values, _ := keyValueList("update", data)

	s := bytes.Buffer{}
	driver := u.getDatabase().Driver
	s.Write(bSQLUpdate)
	s.WriteString(driver.QuoteField(u.table))
	s.Write(bSQLSet)
	s.Write(bytes.Join(keys, BCommaSplit))
	if u.where != nil {
		s.Write(bSQLWhere)
		s.WriteString(u.where.Cond)
		values = append(values, u.where.Args...)
	}
	return s.Bytes(), values
}

func (u *UpdateBuilder) Where(cond string, args ...interface{}) *UpdateBuilder {
	if cond == "" {
		return u
	}
	u.where = &CondParams{cond, args}
	return u
}

func (u *UpdateBuilder) SqlState(data interface{}) *SqlState {
	return NewSqlState(u.parse(data))
}

func (u *UpdateBuilder) Update(data interface{}) (sql.Result, error) {
	return u.getDatabase().ExecPrepare(u.SqlState(data))
}

func (u *UpdateBuilder) UpdateByString(str string) (sql.Result, error) {
	var values []interface{}
	s := bytes.Buffer{}
	driver := u.getDatabase().Driver
	s.Write(bSQLUpdate)
	s.WriteString(driver.QuoteField(u.table))
	s.Write(bSQLSet)
	s.WriteString(str)
	if u.where != nil {
		s.Write(bSQLWhere)
		s.WriteString(u.where.Cond)
		values = u.where.Args
	}

	return u.getDatabase().ExecPrepare(NewSqlState(s.Bytes(), values))
}

func (u *UpdateBuilder) TxUpdate(tx *Tx, data interface{}) (sql.Result, error) {
	return tx.ExecPrepare(u.SqlState(data))
}
