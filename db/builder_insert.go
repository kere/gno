package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"
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

func (ins *InsertBuilder) parseM(rows DataSet) string {
	size := len(rows)
	if size == 0 {
		return ""
	}

	keys, _, _ := keyValueList(ActionInsert, rows[0])
	s := bytes.Buffer{}
	driver := ins.GetDatabase().Driver
	s.Write(insInto)
	s.WriteString(driver.QuoteField(ins.table))
	s.Write(insVal1)
	s.Write(bytes.Join(keys, BCommaSplit))
	s.Write(insVal2)

	// length := len(keys)
	var values []string
	var values2 []interface{}
	var keys2 [][]byte
	var key []byte
	var key2 []byte
	var k int

	for i := 0; i < size; i++ {
		s.Write(insBracketL)
		keys2, values2, _ = keyValueList(ActionInsert, rows[i])
		// 顺序原因，需要重新定位
		values = make([]string, 0)
		for _, key = range keys {
			for k, key2 = range keys2 {
				if string(key) == string(key2) {
					switch values2[k].(type) {
					case time.Time:
						values = append(values, SQuot+(values2[k].(time.Time)).Format(time.RFC1123)+SQuot)
					case string, []byte:
						values = append(values, SQuot+fmt.Sprint(values2[k])+SQuot)
					default:
						values = append(values, fmt.Sprint(values2[k]))
					}
					break
				}
			}
		}
		s.WriteString(strings.Join(values, SCommaSplit))
		s.Write(insBracketR)
		if i < size-1 {
			s.Write(BCommaSplit)
		}
	}

	return s.String()
}

func (ins *InsertBuilder) parse(data interface{}) (string, []interface{}) {
	keys, values, stmts := keyValueList(ActionInsert, data)

	s := bytes.Buffer{}
	driver := ins.GetDatabase().Driver
	s.Write(insInto)
	s.WriteString(driver.QuoteField(ins.table))
	s.Write(insVal1)
	s.Write(bytes.Join(keys, BCommaSplit))
	s.Write(insVal2)
	s.Write(insBracketL)
	s.Write(bytes.Join(stmts, BCommaSplit))
	s.Write(insBracketR)

	return s.String(), values
}

// Insert db
func (ins *InsertBuilder) Insert(data interface{}) (sql.Result, error) {
	cdb := ins.GetDatabase()
	sql, vals := ins.parse(data)
	if ins.isExec {
		return cdb.Exec(sql, vals...)
	}
	return cdb.ExecPrepare(sql, vals...)
}

// InsertM func
func (ins *InsertBuilder) InsertM(rows DataSet) (sql.Result, error) {
	size := 500
	n := len(rows)
	if n <= size {
		sql := ins.parseM(rows)
		return ins.GetDatabase().Exec(sql)
	}

	// pagination
	p := int(math.Ceil(float64(n) / float64(size)))
	var k = 0
	var err error
	var tmp DataSet
	var sqlR sql.Result
	for i := 0; i < p; i++ {
		tmp = DataSet{}
		if i+1 == p {
			//last page
			for k = size * i; k < n; k++ {
				tmp = append(tmp, rows[k])
			}

			sqlR, err = ins.GetDatabase().Exec(ins.parseM(tmp))
			if err != nil {
				return sqlR, err
			}
		} else {
			for k = 0; k < size; k++ {
				tmp = append(tmp, rows[size*i+k])
			}
			sqlR, err = ins.GetDatabase().Exec(ins.parseM(tmp))
			if err != nil {
				return sqlR, err
			}
		}
	}

	return sqlR, nil
}

// TxInsert return sql.Result transation
func (ins *InsertBuilder) TxInsert(tx *Tx, data interface{}) (sql.Result, error) {
	sql, vals := ins.parse(data)
	if ins.isReturnID && tx.GetDatabase().Driver.Name() == "postgres" {
		// s.SetSql(append(s.GetSql(), bPGReturning...))
		sql = sql + sPGReturning

		r, err := tx.Query(sql, vals...)
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

// LastInsertId return int
// func (ins *InsertBuilder) LastInsertId(pkey string) int {
// 	if ins.conn == nil {
// 		ins.GetDatabase().Log.Warn("InsertBuider conn is nil")
// 		return -1
// 	}
//
// 	dataset, err := Query(ins.conn, "SELECT LAST_INSERT_ID() as count")
// 	if err != nil || dataset.IsEmpty() {
// 		ins.GetDatabase().Log.Error(err)
// 		return -1
// 	}
// 	id, err := strconv.ParseInt(dataset[0].String("count"), 10, 64)
// 	if err != nil {
// 		ins.GetDatabase().Log.Error(err)
// 		return -1
// 	}
// 	return int(id)
// }
