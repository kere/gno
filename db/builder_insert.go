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
	isPrepare     bool
}

// NewInsertBuilder func
func NewInsertBuilder(t string) *InsertBuilder {
	return (&InsertBuilder{}).IsPrepare(true).Table(t)
}

// IsPrepare func
func (q *InsertBuilder) IsPrepare(v bool) *InsertBuilder {
	q.isPrepare = v
	return q
}

func (q *InsertBuilder) GetIsPrepare() bool {
	return q.isPrepare
}

func (ins *InsertBuilder) AddExcludeFields(fields ...string) *InsertBuilder {
	if ins.excludeFields == nil {
		ins.excludeFields = make([]string, 0)
	}

	ins.excludeFields = append(ins.excludeFields, fields...)
	return ins
}

func (ins *InsertBuilder) Table(t string) *InsertBuilder {
	ins.table = t
	return ins
}

var b_PG_RETURNING = []byte(" RETURNING id")

func (ins *InsertBuilder) parseM(rows DataSet) []byte {
	size := len(rows)
	if size == 0 {
		return nil
	}

	keys, _, _ := keyValueList("insert", rows[0], ins.excludeFields)
	s := bytes.Buffer{}
	driver := ins.getDatabase().Driver
	s.WriteString("insert into ")
	s.WriteString(driver.QuoteField(ins.table))
	s.WriteString(" (")
	s.Write(bytes.Join(keys, B_CommaSplit))
	s.WriteString(") values ")

	// length := len(keys)
	var values []string
	var values2 []interface{}
	var keys2 [][]byte
	var key []byte
	var key2 []byte
	var k int

	for i := 0; i < size; i++ {
		s.WriteString("(")
		keys2, values2, _ = keyValueList("insert", rows[i], ins.excludeFields)
		// 顺序原因，需要重新定位
		values = make([]string, 0)
		for _, key = range keys {
			for k, key2 = range keys2 {
				if string(key) == string(key2) {
					switch values2[k].(type) {
					case time.Time:
						values = append(values, fmt.Sprintf("'%s'", (values2[k].(time.Time)).Format("2006-01-02 15:04:05")))
					case string, []byte:
						values = append(values, "'"+fmt.Sprint(values2[k])+"'")
					default:
						values = append(values, fmt.Sprint(values2[k]))
					}
					break
				}
			}
		}
		s.WriteString(strings.Join(values, ","))
		s.WriteString(")")
		if i < size-1 {
			s.WriteString(",")
		}
	}

	return s.Bytes()
}

func (ins *InsertBuilder) parse(data interface{}) ([]byte, []interface{}) {
	keys, values, stmts := keyValueList("insert", data, ins.excludeFields)

	s := bytes.Buffer{}
	driver := ins.getDatabase().Driver
	s.WriteString("insert into ")
	s.WriteString(driver.QuoteField(ins.table))
	s.WriteString(" (")
	s.Write(bytes.Join(keys, B_CommaSplit))
	s.WriteString(") values (")
	s.Write(bytes.Join(stmts, B_CommaSplit))
	s.WriteString(")")

	return s.Bytes(), values
}

func (ins *InsertBuilder) SqlState(data interface{}) *SqlState {
	return NewSqlState(ins.parse(data))
}

func (ins *InsertBuilder) Insert(data interface{}) (sql.Result, error) {
	s := ins.SqlState(data)

	cdb := ins.getDatabase()
	if ins.isPrepare {
		return cdb.ExecPrepare(s)
	}

	return cdb.Exec(s)
}

// InsertM
func (ins *InsertBuilder) InsertM(rows DataSet) (sql.Result, error) {
	size := 500
	n := len(rows)
	if n <= size {
		ss := NewSqlState(ins.parseM(rows), nil)
		return ins.getDatabase().Exec(ss)
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

			sqlR, err = ins.getDatabase().Exec(NewSqlState(ins.parseM(tmp), nil))
			if err != nil {
				return sqlR, err
			}
		} else {
			for k = 0; k < size; k++ {
				tmp = append(tmp, rows[size*i+k])
			}
			sqlR, err = ins.getDatabase().Exec(NewSqlState(ins.parseM(tmp), nil))
			if err != nil {
				return sqlR, err
			}
		}
	}

	return sqlR, nil
}

// TxInsert return sql.Result transation
func (ins *InsertBuilder) TxInsert(tx *Tx, data interface{}) (sql.Result, error) {
	s := ins.SqlState(data)
	cdb := ins.getDatabase()
	if cdb.Driver.DriverName() == "postgres" {
		s.SetSql(append(s.GetSql(), b_PG_RETURNING...))

		r, err := tx.Query(s)
		if err != nil {
			return nil, err
		}
		return &insResult{r[0].Int64("id"), 1}, nil
	}
	if ins.isPrepare {
		return tx.ExecPrepare(s)
	}
	return tx.Exec(s)
}

// LastInsertId return int
// func (ins *InsertBuilder) LastInsertId(pkey string) int {
// 	if ins.conn == nil {
// 		ins.getDatabase().Log.Warn("InsertBuider conn is nil")
// 		return -1
// 	}
//
// 	dataset, err := Query(ins.conn, "SELECT LAST_INSERT_ID() as count")
// 	if err != nil || dataset.IsEmpty() {
// 		ins.getDatabase().Log.Error(err)
// 		return -1
// 	}
// 	id, err := strconv.ParseInt(dataset[0].String("count"), 10, 64)
// 	if err != nil {
// 		ins.getDatabase().Log.Error(err)
// 		return -1
// 	}
// 	return int(id)
// }
