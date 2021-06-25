package dba

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/kere/gno/libs/util"
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
	Builder
	excludeFields []string
	IsPrepare     bool
	isReturnID    bool
}

// NewInsert func
func NewInsert(t string) InsertBuilder {
	q := InsertBuilder{IsPrepare: false}
	q.table = t
	return q
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

// Insert db
func (ins *InsertBuilder) Insert(fields []string, row []interface{}) (sql.Result, error) {
	n := len(fields)
	if n != len(row) {
		return nil, fmt.Errorf("insert %s fields.Len() != row.Len() fields: %s row.Len()=%d", ins.table, fields, len(row))
	}
	sqlstr := parseInsert(ins, fields, ins.isReturnID)
	vals := GetRow(n)
	defer PutRow(vals)
	driver := ins.GetDatabase().Driver
	for i := 0; i < n; i++ {
		vals[i] = driver.StoreData(fields[i], row[i])
	}
	if ins.IsPrepare {
		return ins.Exec(sqlstr, vals)
	}
	return ins.ExecPrepare(sqlstr, vals)
}

// Insert0 do PutRow(row) after all
func (ins *InsertBuilder) Insert0(fields []string, row []interface{}) (sql.Result, error) {
	n := len(fields)
	if n != len(row) {
		return nil, fmt.Errorf("insert0 %s fields.Len() != row.Len() fields: %s row.Len()=%d", ins.table, fields, len(row))
	}
	sqlstr := parseInsert(ins, fields, ins.isReturnID)
	defer PutRow(row)
	driver := ins.GetDatabase().Driver
	for i := 0; i < n; i++ {
		row[i] = driver.StoreData(fields[i], row[i])
	}
	if ins.IsPrepare {
		return ins.Exec(sqlstr, row)
	}
	return ins.ExecPrepare(sqlstr, row)
}

// InsertM func
func (ins *InsertBuilder) InsertM(dat *DataSet) (sql.Result, error) {
	sqlstr, vals := parseInsertMP(ins, dat)
	defer PutColumn(vals)
	if ins.IsPrepare {
		return ins.Exec(sqlstr, vals)
	}
	return ins.ExecPrepare(sqlstr, vals)
}

// InsertM0 do PutDataSet after all
func (ins *InsertBuilder) InsertM0(dat *DataSet) (sql.Result, error) {
	sqlstr, vals := parseInsertMP(ins, dat)
	PutDataSet(dat)
	defer PutColumn(vals)
	if ins.IsPrepare {
		return ins.Exec(sqlstr, vals)
	}
	return ins.ExecPrepare(sqlstr, vals)
}

func parseInsertMP(ins *InsertBuilder, dataset *DataSet) (string, []interface{}) {
	if dataset.Len() == 0 {
		return "", nil
	}

	keys := dataset.Fields
	buf := bytebufferpool.Get()

	driver := ins.GetDatabase().Driver

	buf.Write(bInsSQL)
	driver.WriteQuoteIdentifier(buf, ins.table)
	buf.Write(bInsBracketL)
	n := len(keys)
	driver.WriteQuoteIdentifier(buf, keys[0])
	for i := 1; i < n; i++ {
		buf.WriteByte(',')
		driver.WriteQuoteIdentifier(buf, keys[i])
	}
	buf.Write(bInsBracketR)

	values := writeInsertMP(buf, dataset)

	str := buf.String()
	bytebufferpool.Put(buf)

	return str, values
}

func parseInsert(ins *InsertBuilder, fields []string, hasReturnID bool) string {
	l := len(fields)
	stmts := make([][]byte, l)
	for i := 0; i < l; i++ {
		v := []byte(fmt.Sprint(i + 1))
		stmts[i] = append([]byte{'$'}, v...)
	}

	s := bytes.Buffer{}
	driver := ins.GetDatabase().Driver
	s.Write(bInsSQL)
	driver.WriteQuoteIdentifier(&s, ins.table)
	s.Write(bInsBracketL)
	s.Write(util.Str2Bytes(fields[0]))
	for i := 1; i < l; i++ {
		s.WriteRune(',')
		s.Write(util.Str2Bytes(fields[i]))
	}
	s.Write(bInsBracketR)
	s.WriteByte('(')
	s.Write(bytes.Join(stmts, BComma))
	s.WriteByte(')')

	if hasReturnID {
		s.Write(bPGReturning)
	}

	return s.String()
}

func writeInsertMP(buf *bytebufferpool.ByteBuffer, dataset *DataSet) []interface{} {
	l := dataset.Len()
	n := len(dataset.Fields)
	seq := 1

	database := Current()
	// values := make([]interface{}, 0, l*n+10)
	values := GetColumn(0)
	cols := dataset.Columns

	for i := 0; i < l; i++ {
		buf.WriteByte('(')

		for k := 0; k < n; k++ {
			key := dataset.Fields[k]
			v := cols[k][i]
			val := database.Driver.StoreData(key, v)

			values = append(values, val)

			if seq == 1 {
				buf.WriteString("$1")
			} else {
				buf.WriteString(fmt.Sprint(SDoller, seq))
			}
			if k < n-1 {
				buf.WriteByte(',')
			}
			seq++
		}

		buf.WriteByte(')')
		if i < l-1 {
			buf.WriteByte(',')
		}

	}
	buf.WriteByte(';')

	return values
}
