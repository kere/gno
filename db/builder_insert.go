package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"

	"github.com/valyala/bytebufferpool"
)

type insResult struct {
	id           int64
	rowsAffected int64
}

func (i *insResult) LastInsertId() (int64, error) {
	return i.id, nil
}
func (i *insResult) RowsAffected() (int64, error) {
	return i.rowsAffected, nil
}

// InsertBuilder class
type InsertBuilder struct {
	table string
	builder
	excludeFields []string
	isExec        bool
	isReturnID    bool
}

// NewInsert func
func NewInsert(t string) *InsertBuilder {
	return &InsertBuilder{table: t, isExec: false}
}

// SetIsPrepare func
func (ins *InsertBuilder) SetIsPrepare(v bool) *InsertBuilder {
	ins.isExec = !v
	return ins
}

// GetIsPrepare get
func (ins *InsertBuilder) GetIsPrepare() bool {
	return !ins.isExec
}

// ReturnID func
func (ins *InsertBuilder) ReturnID() *InsertBuilder {
	ins.isReturnID = true
	return ins
}

// AddSkipFields skip fields
func (ins *InsertBuilder) AddSkipFields(fields ...string) *InsertBuilder {
	if ins.excludeFields == nil {
		ins.excludeFields = make([]string, 0)
	}

	ins.excludeFields = append(ins.excludeFields, fields...)
	return ins
}

// Table string
func (ins *InsertBuilder) Table(t string) *InsertBuilder {
	ins.table = t
	return ins
}

func parseInsertM(ins *InsertBuilder, rows MapRows) (string, []interface{}) {
	l := len(rows)
	if l == 0 {
		return "", nil
	}

	bkeys, strkeys := sqlInsertKeysByMapRow(rows[0])

	// buf := bytePool.Get()
	buf := bytebufferpool.Get()

	driver := ins.GetDatabase().Driver
	buf.Write(bInsSQL)
	buf.Write(driver.QuoteIdentifierB(ins.table))
	buf.Write(bInsBracketL)
	buf.Write(bkeys)
	buf.Write(bInsBracketR)

	values := writeInsertMByMapRow(buf, strkeys, rows)

	str := buf.String()
	// bytePool.Put(buf)
	bytebufferpool.Put(buf)

	return str, values
}

func parseInsertM2(ins *InsertBuilder, dataset DataSet) (string, []interface{}) {
	if dataset.Len() == 0 {
		return "", nil
	}

	keys := dataset.Fields

	// buf := bytes.NewBuffer(nil)
	// buf := bytePool.Get()
	buf := bytebufferpool.Get()

	database := ins.GetDatabase()
	driver := database.Driver

	buf.Write(bInsSQL)
	buf.Write(driver.QuoteIdentifierB(ins.table))
	buf.Write(bInsBracketL)
	n := len(keys)
	for i := 0; i < n; i++ {
		buf.Write(database.Driver.QuoteIdentifierB(keys[i]))
		if i < n-1 {
			buf.WriteByte(',')
		}
	}
	buf.Write(bInsBracketR)

	values := writeInsertMByDataSet(buf, &dataset)

	str := buf.String()
	bytebufferpool.Put(buf)

	return str, values
}

func parseInsert(ins *InsertBuilder, row MapRow, hasReturnID bool) (string, []interface{}) {
	// keys, values, stmts := keyValueList(ActionInsert, row)
	keys, values := sqlInsertParamsByMapRow(row)
	l := len(values)
	stmts := make([][]byte, l)
	for i := 0; i < l; i++ {
		v := []byte(fmt.Sprint(i + 1))
		stmts[i] = append([]byte{'$'}, v...)
	}

	s := bytes.Buffer{}
	driver := ins.GetDatabase().Driver
	s.Write(bInsSQL)
	s.Write(driver.QuoteIdentifierB(ins.table))
	s.Write(bInsBracketL)
	s.Write(keys)
	s.Write(bInsBracketR)
	s.WriteByte('(')
	s.Write(bytes.Join(stmts, BComma))
	s.WriteByte(')')

	if hasReturnID {
		s.Write(bPGReturning)
	}

	return s.String(), values
}

// Insert db
func (ins *InsertBuilder) Insert(row MapRow) (sql.Result, error) {
	cdb := ins.GetDatabase()
	sql, vals := parseInsert(ins, row, false)
	if ins.isExec {
		return cdb.Exec(sql, vals...)
	}
	return cdb.ExecPrepare(sql, vals...)
}

// InsertMN every n
func (ins *InsertBuilder) InsertMN(rows interface{}, step int) error {
	var isMapRows bool
	var l int
	var drows MapRows
	var dataset DataSet

	switch rows.(type) {
	case MapRows:
		isMapRows = true
		drows = rows.(MapRows)
		l = drows.Len()
	case DataSet:
		isMapRows = false
		dataset = rows.(DataSet)
		l = dataset.Len()
	}

	count := int(math.Ceil(float64(l) / float64(step)))
	for i := 0; i < count; i++ {
		b := i * step
		e := b + step
		if e > l {
			e = l
		}
		if isMapRows {
			_, err := ins.InsertM(drows[b:e])
			if err != nil {
				return err
			}
		} else {
			_, err := ins.InsertM(dataset.RangeI(b, e))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// InsertM func
func (ins *InsertBuilder) InsertM(rows interface{}) (sql.Result, error) {
	var sqlstr string
	var vals []interface{}

	switch rows.(type) {
	case MapRows:
		sqlstr, vals = parseInsertM(ins, rows.(MapRows))
	case DataSet:
		sqlstr, vals = parseInsertM2(ins, rows.(DataSet))
	default:
		return nil, ErrType
	}

	cdb := ins.GetDatabase()

	if ins.isExec {
		return cdb.Exec(sqlstr, vals...)
	}
	return cdb.ExecPrepare(sqlstr, vals...)
}

// TxInsert return sql.Result transation
func (ins *InsertBuilder) TxInsert(tx *Tx, data MapRow) (sql.Result, error) {
	hasReturnID := ins.isReturnID && tx.GetDatabase().Driver.Name() == "postgres"
	sql, vals := parseInsert(ins, data, hasReturnID)
	if hasReturnID {
		r, err := tx.QueryRows(sql, vals...)
		if err != nil {
			return nil, err
		}
		return &insResult{r[0].Int64("id"), 1}, nil
	}
	if ins.isExec {
		return tx.Exec(sql, vals...)
	}
	return tx.ExecPrepare(sql, vals...)
}

// TxInsertM return sql.Result transation
func (ins *InsertBuilder) TxInsertM(tx *Tx, rows interface{}) (sql.Result, error) {
	var sqlstr string
	var vals []interface{}

	switch rows.(type) {
	case MapRows:
		sqlstr, vals = parseInsertM(ins, rows.(MapRows))
	case DataSet:
		sqlstr, vals = parseInsertM2(ins, rows.(DataSet))
	default:
		return nil, ErrType
	}

	if ins.isExec {
		return tx.Exec(sqlstr, vals...)
	}
	return tx.ExecPrepare(sqlstr, vals...)
}

// TxInsertMN return sql.Result transation
func (ins *InsertBuilder) TxInsertMN(tx *Tx, rows interface{}, step int) error {
	var isMapRows bool
	var l int
	var drows MapRows
	var dataset DataSet

	switch rows.(type) {
	case MapRows:
		isMapRows = true
		drows = rows.(MapRows)
		l = drows.Len()
	case DataSet:
		isMapRows = false
		dataset = rows.(DataSet)
		l = dataset.Len()
	}

	count := int(math.Ceil(float64(l) / float64(step)))
	for i := 0; i < count; i++ {
		b := i * step
		e := b + step
		if e > l {
			e = l
		}
		if isMapRows {
			_, err := ins.TxInsertM(tx, drows[b:e])
			if err != nil {
				return err
			}
		} else {
			_, err := ins.TxInsertM(tx, dataset.RangeI(b, e))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
