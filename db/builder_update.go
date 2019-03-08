package db

import (
	"bytes"
	"database/sql"
)

// UpdateBuilder class
type UpdateBuilder struct {
	table string
	where string
	args  []interface{}
	builder
}

// NewUpdate func
func NewUpdate(t string) *UpdateBuilder {
	return &UpdateBuilder{table: t}
}

// Table string
func (u *UpdateBuilder) Table(t string) *UpdateBuilder {
	u.table = t
	return u
}

func (u *UpdateBuilder) parse(data interface{}) (string, []interface{}) {
	keys, values, _ := keyValueList(ActionUpdate, data)

	s := bytes.Buffer{}
	driver := u.GetDatabase().Driver
	s.Write(bSQLUpdate)
	s.WriteString(driver.QuoteField(u.table))
	s.Write(bSQLSet)
	s.Write(bytes.Join(keys, BCommaSplit))
	if u.where != "" {
		s.Write(bSQLWhere)
		s.WriteString(driver.Adapt(u.where, len(values)))
		values = append(values, u.args...)
	}
	return s.String(), values
}

// Where sql
func (u *UpdateBuilder) Where(cond string, args ...interface{}) *UpdateBuilder {
	if cond == "" {
		return u
	}
	u.where = cond
	u.args = args
	return u
}

// Update db
func (u *UpdateBuilder) Update(data interface{}) (sql.Result, error) {
	sql, vals := u.parse(data)
	return u.GetDatabase().ExecPrepare(sql, vals...)
}

// // UpdateByString by string
// func (u *UpdateBuilder) UpdateByString(str string) (sql.Result, error) {
// 	var values []interface{}
// 	s := bytes.Buffer{}
// 	driver := u.GetDatabase().Driver
// 	s.Write(bSQLUpdate)
// 	s.WriteString(driver.QuoteField(u.table))
// 	s.Write(bSQLSet)
// 	s.WriteString(str)
// 	if u.where != "" {
// 		s.Write(bSQLWhere)
// 		s.WriteString(u.where)
// 		values = u.args
// 	}
//
// 	return u.GetDatabase().ExecPrepare(s.String(), values...)
// }

// // TxUpdateByString trunsaction
// func (u *UpdateBuilder) TxUpdateByString(tx *Tx, str string) (sql.Result, error) {
// 	var values []interface{}
// 	s := bytes.Buffer{}
// 	driver := u.GetDatabase().Driver
// 	s.Write(bSQLUpdate)
// 	s.WriteString(driver.QuoteField(u.table))
// 	s.Write(bSQLSet)
// 	s.WriteString(str)
// 	if u.where != "" {
// 		s.Write(bSQLWhere)
// 		s.WriteString(u.where)
// 		values = u.args
// 	}
//
// 	return tx.ExecPrepare(s.String(), values...)
// }

// TxUpdate trunsaction
func (u *UpdateBuilder) TxUpdate(tx *Tx, data interface{}) (sql.Result, error) {
	sql, vals := u.parse(data)
	return tx.ExecPrepare(sql, vals...)
}
