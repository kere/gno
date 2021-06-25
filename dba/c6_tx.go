package dba

import (
	"database/sql"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/myerr"
)

type Tx struct {
	tx        *sql.Tx
	database  *Database
	LastError error
}

// BeginTx tx
func BeginTx() (Tx, error) {
	var err error
	t := Tx{database: Current()}
	t.tx, err = t.database.DB().Begin()
	if err != nil {
		return t, err
	}
	return t, nil
}

// End func
func (t *Tx) End() error {
	t.LastError = nil
	return t.Commit()
}

// Commit func
func (t *Tx) Commit() error {
	err := t.tx.Commit()
	if err != nil {
		log.App.Alert(err)
	}
	return err
}

// DoError tx
func (t *Tx) DoError(err error) bool {
	if err != nil {
		myerr.New(err).Log().Stack()
		err2 := t.tx.Rollback()
		if err2 != nil {
			return true
		}
		t.LastError = err
		return true
	}
	return false
}

// Rollback err
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

// NewBuilder
func (t *Tx) NewBuilder(table string) Builder {
	b := Builder{table: table}
	b.database = t.database
	b.isTx = true
	b.tx = t.tx
	return b
}

// NewQuery
func (t *Tx) NewQuery(table string) QueryBuilder {
	q := newQuery(table)
	q.database = t.database
	q.isTx = true
	q.tx = t.tx
	return q
}

// NewInsert
func (t *Tx) NewInsert(table string) InsertBuilder {
	ins := newInsert(table)
	ins.database = t.database
	ins.isTx = true
	ins.tx = t.tx
	return ins
}

// NewUpdate
func (t *Tx) NewUpdate(table string) UpdateBuilder {
	u := newUpdate(table)
	u.database = t.database
	u.isTx = true
	u.tx = t.tx
	return u
}

// NewDelete
func (t *Tx) NewDelete(table string) DeleteBuilder {
	del := newDelete(table)
	del.database = t.database
	del.isTx = true
	del.tx = t.tx
	return del
}

// NewExists
func (t *Tx) NewExists(table string) ExistsBuilder {
	e := newExists(table)
	e.database = t.database
	e.isTx = true
	e.tx = t.tx
	return e
}
