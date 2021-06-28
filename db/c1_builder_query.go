package db

import (
	"fmt"
	"strconv"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
)

// QueryBuilder class
type QueryBuilder struct {
	Builder

	where string
	args  []interface{}

	fields   string
	leftJoin string
	order    string
	limit    int
	offset   int
}

// NewQuery new
func NewQuery(t string) QueryBuilder {
	q := QueryBuilder{}
	q.table = t
	return q
}

// Table name
func (q *QueryBuilder) Table(t string) *QueryBuilder {
	q.table = t
	return q
}

// SetPrepare prepare sql
func (q *QueryBuilder) Prepare(v bool) *QueryBuilder {
	q.isPrepare = v
	return q
}

// Select fields
func (q *QueryBuilder) Select(s string) *QueryBuilder {
	q.fields = s
	return q
}

// LeftJoin sql
func (q *QueryBuilder) LeftJoin(s string) *QueryBuilder {
	q.leftJoin = s
	return q
}

// Where sql
func (q *QueryBuilder) Where(s string, args ...interface{}) *QueryBuilder {
	q.where = s
	q.args = args
	return q
}

// GetWhere sql
func (q *QueryBuilder) GetWhere() (string, []interface{}) {
	return q.where, q.args
}

// Order sql
func (q *QueryBuilder) Order(s string) *QueryBuilder {
	q.order = s
	return q
}

// Page pagation
func (q *QueryBuilder) Page(page int, pageSize int) *QueryBuilder {
	q.offset = (page - 1) * pageSize
	q.limit = pageSize
	return q
}

// Limit sql
func (q *QueryBuilder) Limit(n int) *QueryBuilder {
	q.limit = n
	return q
}

// Offset sql
func (q *QueryBuilder) Offset(n int) *QueryBuilder {
	q.offset = n
	return q
}

// Parse sql
func (q *QueryBuilder) Parse() string {
	buf := bytebufferpool.Get()

	buf.Write(bSQLSelect)

	setQueryFields(q, buf)

	buf.Write(bSQLFrom)
	buf.WriteString(q.table)

	if q.leftJoin != "" {
		buf.Write(bSQLLeftJoin)
		buf.WriteString(q.leftJoin)
	}

	if q.where != "" {
		buf.Write(bSQLWhere)
		buf.WriteString(q.where)
	}

	if q.order != "" {
		buf.Write(bSQLOrder)
		buf.WriteString(q.order)
	}

	limit := q.limit

	if limit > 0 && q.offset > 0 {
		buf.Write(bSQLLimit)
		buf.WriteString(fmt.Sprint(limit))
		buf.Write(bSQLOffset)
		buf.WriteString(fmt.Sprint(q.offset))
	} else if limit > 0 {
		buf.Write(bSQLLimit)
		buf.WriteString(strconv.FormatInt(int64(limit), 10))
	} else if q.offset > 0 {
		buf.Write(bSQLOffset)
		buf.WriteString(strconv.FormatInt(int64(q.offset), 10))
	}
	str := buf.String()
	// bytePool.Put(buf)
	bytebufferpool.Put(buf)

	return str
}

// QueryP return DataSet
func (q *QueryBuilder) QueryP() (DataSet, error) {
	sqlstr := q.Parse()
	return q.cQuery(true, sqlstr, q.args)
}

// Query return DataSet
func (q *QueryBuilder) Query() (DataSet, error) {
	sqlstr := q.Parse()
	return q.cQuery(false, sqlstr, q.args)
}

// QueryOne limit=1
func (q *QueryBuilder) QueryOne() (DBRow, error) {
	return q.queryOne(false)
}

// QueryOneP limit=1
func (q *QueryBuilder) QueryOneP() (DBRow, error) {
	return q.queryOne(true)
}

// QueryOne limit=1
func (q *QueryBuilder) queryOne(isPool bool) (DBRow, error) {
	limit := q.limit
	q.limit = 1
	sqlstr := q.Parse()

	dataset, err := q.cQuery(true, sqlstr, q.args)
	defer PutDataSet(&dataset)
	q.limit = limit
	if err != nil {
		return DBRow{}, err
	}

	if dataset.Len() == 0 {
		return DBRow{Fields: dataset.Fields}, err
	}
	if isPool {
		return DBRow{Values: dataset.RowAtP(0), Fields: dataset.Fields}, nil
	}
	return DBRow{Values: dataset.RowAt(0), Fields: dataset.Fields}, nil
}

func setQueryFields(q *QueryBuilder, buf *bytebufferpool.ByteBuffer) {
	fields := util.BStarKey
	if q.fields != "" {
		fields = []byte(q.fields)
	}

	buf.Write(fields)
}
