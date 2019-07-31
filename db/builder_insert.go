package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"
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

var (
	insInto = []byte("insert into ")
	insVal1 = []byte(" (")
	insVal2 = []byte(") values ")
	// insBracketL = []byte("(")
	// insBracketR = []byte(")")
)

func parseInsertM(ins *InsertBuilder, rows MapRows) (string, []interface{}) {
	l := len(rows)
	if l == 0 {
		return "", nil
	}

	bkeys, strkeys := sqlInsertKeysByMapRow(rows[0])

	buf := bytes.NewBuffer(nil)
	driver := ins.GetDatabase().Driver
	buf.Write(insInto)
	buf.WriteString(driver.QuoteField(ins.table))
	buf.Write(insVal1)
	buf.Write(bkeys)
	buf.Write(insVal2)

	values := writeInsertMByMapRow(buf, strkeys, rows)

	return buf.String(), values
}

func parseInsertM2(ins *InsertBuilder, dataset *DataSet) (string, []interface{}) {
	if dataset.Len() == 0 {
		return "", nil
	}

	keys := dataset.Fields

	buf := bytes.NewBuffer(nil)
	database := ins.GetDatabase()
	driver := database.Driver

	buf.Write(insInto)
	buf.WriteString(driver.QuoteField(ins.table))
	buf.Write(insVal1)
	n := len(keys)
	for i := 0; i < n; i++ {
		buf.Write(database.Driver.QuoteFieldB(keys[i]))
		if i < n-1 {
			buf.WriteByte(',')
		}
	}
	buf.Write(insVal2)

	values := writeInsertMByDataSet(buf, dataset)

	return buf.String(), values
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
	s.Write(insInto)
	s.WriteString(driver.QuoteField(ins.table))
	s.Write(insVal1)
	s.Write(keys)
	s.Write(insVal2)
	s.WriteByte('(')
	s.Write(bytes.Join(stmts, BCommaSplit))
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
func (ins *InsertBuilder) InsertMN(rows MapRows, step int) error {
	l := rows.Len()
	count := int(math.Ceil(float64(l) / float64(step)))
	for i := 0; i < count; i++ {
		b := i * step
		e := b + step
		if e > l {
			e = l
		}
		_, err := ins.InsertM(rows[b:e])
		if err != nil {
			return err
		}
	}
	return nil
}

// InsertM func
func (ins *InsertBuilder) InsertM(rows MapRows) (sql.Result, error) {
	str, vals := parseInsertM(ins, rows)
	cdb := ins.GetDatabase()

	if ins.isExec {
		return cdb.Exec(str, vals...)
	}
	return cdb.ExecPrepare(str, vals...)
}

// InsertMN2 every n
func (ins *InsertBuilder) InsertMN2(dataset DataSet, step int) error {
	l := dataset.Len()
	count := int(math.Ceil(float64(l) / float64(step)))
	for i := 0; i < count; i++ {
		b := i * step
		e := b + step
		if e > l {
			e = l
		}
		_, err := ins.InsertM2(dataset.RangeI(b, e))
		if err != nil {
			return err
		}
	}
	return nil
}

// InsertM2 func
func (ins *InsertBuilder) InsertM2(dataset DataSet) (sql.Result, error) {
	str, vals := parseInsertM2(ins, &dataset)
	cdb := ins.GetDatabase()

	if ins.isExec {
		return cdb.Exec(str, vals...)
	}
	return cdb.ExecPrepare(str, vals...)
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
func (ins *InsertBuilder) TxInsertM(tx *Tx, rows MapRows) (sql.Result, error) {
	str, vals := parseInsertM(ins, rows)
	if ins.isExec {
		return tx.Exec(str, vals...)
	}
	return tx.ExecPrepare(str, vals...)
}

// TxInsertMN return sql.Result transation
func (ins *InsertBuilder) TxInsertMN(tx *Tx, rows MapRows, step int) error {
	l := rows.Len()
	count := int(math.Ceil(float64(l) / float64(step)))
	for i := 0; i < count; i++ {
		b := i * step
		e := b + step
		if e > l {
			e = l
		}
		_, err := ins.TxInsertM(tx, rows[b:e])
		if err != nil {
			return err
		}
	}
	return nil
}

// TxInsertM2 return sql.Result transation
func (ins *InsertBuilder) TxInsertM2(tx *Tx, dataset DataSet) (sql.Result, error) {
	str, vals := parseInsertM2(ins, &dataset)
	if ins.isExec {
		return tx.Exec(str, vals...)
	}
	return tx.ExecPrepare(str, vals...)
}

// TxInsertMN2 return sql.Result transation
func (ins *InsertBuilder) TxInsertMN2(tx *Tx, dataset DataSet, step int) error {
	l := dataset.Len()
	count := int(math.Ceil(float64(l) / float64(step)))
	for i := 0; i < count; i++ {
		b := i * step
		e := b + step
		if e > l {
			e = l
		}
		_, err := ins.TxInsertM2(tx, dataset.RangeI(b, e))
		if err != nil {
			return err
		}
	}
	return nil
}
