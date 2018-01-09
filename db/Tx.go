package db

import (
	"database/sql"
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

func (t *Tx) Begin() *Tx {
	t.GetTx()
	return t
}

func (t *Tx) GetTx() *sql.Tx {
	if t.tx != nil {
		return t.tx
	}

	t.conn = Current().Connection.Connect()
	tx, err := t.conn.Begin()
	if err != nil {
		Current().Log.Error(err)
		return nil
	}
	t.tx = tx
	t.IsError = false
	return t.tx
}

func (t *Tx) FindOne(cls IVO, item *SqlState) (IVO, error) {
	r, err := t.Find(cls, item)
	if err != nil {
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
	if err != nil {
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
	if err != nil {
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

func (t *Tx) End() bool {
	defer t.close()
	return t.Commit()
}

func (t *Tx) Commit() bool {
	if t.IsError {
		return false
	}
	if err := t.tx.Commit(); err != nil {
		panic(err)
	}

	return true
}

func (t *Tx) DoError(err error) bool {
	if err != nil {
		Current().Log.Error(err)
		t.tx.Rollback()
		t.IsError = true
		t.LastError = err
		return true
	}
	return false
}
