package db

import (
	"database/sql"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/myerr"
)

type Tx struct {
	builder
	conn      *sql.DB
	tx        *sql.Tx
	IsError   bool
	LastError error
}

func NewTx() *Tx {
	return &Tx{}
}

func (t *Tx) Begin() error {
	// if t.tx != nil {
	// 	return t
	// }

	t.conn = Current().Connection.Connect()
	tx, err := t.conn.Begin()
	if err != nil {
		t.IsError = true
		log.App.Alert(err).Stack()
		return err
	}
	t.tx = tx
	t.IsError = false
	return nil
}

func (t *Tx) FindOne(cls IVO, item *SqlState) (IVO, error) {
	r, err := t.Find(cls, item)
	if t.DoError(err) {
		return nil, err
	}

	if len(r) > 0 {
		return r[0], nil
	} else {
		return nil, nil
	}
}

func (t *Tx) QueryOne(item *SqlState) (DataRow, error) {
	r, err := t.Query(item)
	if t.DoError(err) {
		return nil, err
	}

	if len(r) > 0 {
		return r[0], nil
	} else {
		return nil, nil
	}
}

func (t *Tx) Exists(item *SqlState) (bool, error) {
	r, err := t.Query(item)
	if t.DoError(err) {
		return false, err
	}
	return len(r) > 0, nil
}

func (t *Tx) Query(s *SqlState) (DataSet, error) {
	if t.IsError {
		return nil, nil
	}

	// sqlstr := string(t.database.AdaptSql(item.Sql))
	database := Current()
	bsqlstr := s.GetSql()

	database.Log.Sql(bsqlstr, s.GetArgs())

	st, err := t.tx.Prepare(string(bsqlstr))
	if t.DoError(err) {
		return nil, err
	}

	rows, err := st.Query(s.GetArgs()...)
	if t.DoError(err) {
		return nil, err
	}

	defer rows.Close()

	dataset, err := ScanRows(rows)
	if t.DoError(err) {
		return nil, err
	}

	return dataset, nil
}

func (t *Tx) Find(cls IVO, item *SqlState) (VODataSet, error) {
	if t.IsError {
		return nil, nil
	}

	// sqlstr := string(t.database.AdaptSql(item.Sql))
	database := Current()
	database.Log.Sql(item.GetSql(), item.GetArgs())

	dataset, err := t.Query(item)
	if t.DoError(err) {
		return nil, err
	}

	return NewStructConvert(cls).DataSet2Struct(dataset)
}

func (t *Tx) Exec(s *SqlState) (sql.Result, error) {
	if t.IsError {
		return nil, nil
	}
	// sqlstr := string(t.database.AdaptSql(item.Sql))
	database := Current()
	bsqlstr := s.GetSql()

	database.Log.Sql(bsqlstr, s.GetArgs())

	r, err := t.tx.Exec(string(bsqlstr), s.GetArgs()...)
	if t.DoError(err) {
		return nil, err
	}
	return r, nil
}

func (t *Tx) LastInsertId(table, pkey string) int64 {
	if t.IsError {
		return -1
	}
	r := t.tx.QueryRow(Current().Driver.LastInsertID(table, pkey))

	var count int64
	err := r.Scan(&count)
	if t.DoError(err) {
		return -1
	}
	return count
}

func (t *Tx) ExecPrepare(s *SqlState) (sql.Result, error) {
	if t.IsError {
		return nil, nil
	}
	// sqlstr := string(t.database.AdaptSql(item.Sql))
	database := Current()
	bsqlstr := s.GetSql()

	database.Log.Sql(bsqlstr, s.GetArgs())

	sqlstr := string(bsqlstr)
	st, err := t.tx.Prepare(sqlstr)
	if t.DoError(err) {
		return nil, err
	}
	r, err := st.Exec(s.GetArgs()...)
	if t.DoError(err) {
		return nil, err
	}
	return r, nil
}

func (t *Tx) close() error {
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

func (t *Tx) End() error {
	err := t.Commit()
	if err != nil {
		return err
	}
	return t.close()
}

func (t *Tx) Commit() error {
	if t.IsError {
		return nil
	}
	err := t.tx.Commit()
	if err != nil {
		log.App.Alert(err)
	}
	return err
}

// DoError tx
func (t *Tx) DoError(err error) bool {
	if t.IsError {
		return true
	}

	if err != nil {
		myerr.New(err).Log().Stack()
		err2 := t.tx.Rollback()
		if err2 != nil {
			log.App.Alert(err2, "rollback faield")
			return true
		}
		t.IsError = true
		t.LastError = err
		return true
	}
	return false
}

// Rollback err
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
