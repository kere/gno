package db

import (
	"bytes"
	"database/sql"
	"fmt"
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
	insInto     = []byte("insert into ")
	insVal1     = []byte(" (")
	insVal2     = []byte(") values ")
	insBracketL = []byte("(")
	insBracketR = []byte(")")
)

// func parseInsertM(ins *InsertBuilder, rows MapRows) string {
// 	l := len(rows)
// 	if l == 0 {
// 		return ""
// 	}
//
// 	// keys, _, _ := keyValueList(ActionInsert, rows[0])
// 	bkeys, strkeys := sqlInsertKeysByMapRow(rows[0])
// 	s := bytes.Buffer{}
// 	driver := ins.GetDatabase().Driver
// 	s.Write(insInto)
// 	s.WriteString(driver.QuoteField(ins.table))
// 	s.Write(insVal1)
// 	s.Write(bkeys)
// 	s.Write(insVal2)
//
//
// 	for i := 0; i < l; i++ {
// 		s.Write(insBracketL)
// 		// keys2, values2, _ = keyValueList(ActionInsert, rows[i])
// 		valstr := sqlInsertStringByMapRow(strkeys, rows[i])
// 		s.WriteString(valstr)
// 		s.Write(insBracketR)
// 		if i < l-1 {
// 			s.Write(BCommaSplit)
// 		}
// 	}
//
// 	return s.String()
// }

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
	s.Write(insBracketL)
	s.Write(bytes.Join(stmts, BCommaSplit))
	s.Write(insBracketR)

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

// // InsertMN every n
// func (ins *InsertBuilder) InsertMN(rows MapRows, step int) error {
// 	l := rows.Len()
// 	count := int(math.Ceil(float64(l) / float64(step)))
// 	for i := 0; i < count; i++ {
// 		b := i * step
// 		e := b + step
// 		if e > l {
// 			e = l
// 		}
// 		_, err := ins.InsertM(rows[b:e])
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // InsertM func
// func (ins *InsertBuilder) InsertM(rows MapRows) (int, error) {
// 	l := len(rows)
// 	step := 500
//
// 	if l <= step {
// 		sql := parseInsertM(ins, rows)
// 		_, err := ins.GetDatabase().Exec(sql)
// 		return l, err
// 	}
//
// 	n := l/step + 1
// 	count := 0
// 	for i := 0; i < n; i++ {
// 		b := i * step
// 		if b >= l {
// 			break
// 		}
// 		e := b + step
// 		if e > l {
// 			e = l
// 		}
// 		_, err := ins.GetDatabase().Exec(parseInsertM(ins, rows[b:e]))
// 		if err != nil {
// 			return count, err
// 		}
// 		count += e - b
// 	}
// 	return count, nil
// }

// // InsertM func
// func (ins *InsertBuilder) InsertM(rows MapRows) (sql.Result, error) {
// 	size := 500
// 	n := len(rows)
// 	if n <= size {
// 		sql := parseInsertM(ins, rows)
// 		return ins.GetDatabase().Exec(sql)
// 	}
//
// 	// pagination
// 	p := int(math.Ceil(float64(n) / float64(size)))
// 	var k int
// 	var err error
// 	var tmp MapRows
// 	var sqlR sql.Result
// 	for i := 0; i < p; i++ {
// 		tmp = MapRows{}
// 		if i+1 == p {
// 			//last page
// 			for k = size * i; k < n; k++ {
// 				tmp = append(tmp, rows[k])
// 			}
//
// 			sqlR, err = ins.GetDatabase().Exec(parseInsertM(ins, tmp))
// 			if err != nil {
// 				return sqlR, err
// 			}
// 		} else {
// 			for k = 0; k < size; k++ {
// 				tmp = append(tmp, rows[size*i+k])
// 			}
// 			sqlR, err = ins.GetDatabase().Exec(parseInsertM(ins, tmp))
// 			if err != nil {
// 				return sqlR, err
// 			}
// 		}
// 	}
//
// 	return sqlR, nil
// }

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
