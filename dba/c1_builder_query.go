package dba

import (
	"fmt"
	"strconv"

	"github.com/valyala/bytebufferpool"
)

// QueryBuilder class
type QueryBuilder struct {
	Builder

	fields    string
	leftJoin  string
	where     string
	args      []interface{}
	order     string
	limit     int
	offset    int
	isPrepare bool
}

// NewQuery new
func NewQuery(t string) QueryBuilder {
	q := QueryBuilder{isPrepare: true}
	q.table = t
	return q
}

// SetPrepare prepare sql
func (q *QueryBuilder) SetPrepare(v bool) *QueryBuilder {
	q.isPrepare = v
	return q
}

// GetPrepare get
func (q *QueryBuilder) GetPrepare() bool {
	return q.isPrepare
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
	if s == "" {
		return q
	}

	q.where = s
	q.args = args
	return q
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
func (q *QueryBuilder) Parse() (string, []interface{}) {
	// s := bytes.Buffer{}
	// buf := bytePool.Get()
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

	return str, q.args
}

// QueryP return DataSet
func (q *QueryBuilder) QueryP() (DataSet, error) {
	sqlstr, args := q.Parse()
	return q.cQuery(true, q.isPrepare, sqlstr, args)
}

// Query return DataSet
func (q *QueryBuilder) Query() (DataSet, error) {
	sqlstr, args := q.Parse()
	return q.cQuery(false, q.isPrepare, sqlstr, args)
}

// QueryOne limit=1
func (q *QueryBuilder) QueryOne() ([]interface{}, []string, error) {
	return q.queryOne(false)
}

// QueryOneP limit=1
func (q *QueryBuilder) QueryOneP() ([]interface{}, []string, error) {
	return q.queryOne(true)
}

// QueryOne limit=1
func (q *QueryBuilder) queryOne(isPool bool) ([]interface{}, []string, error) {
	limit := q.limit
	q.limit = 1
	dataset, err := q.QueryP()
	defer PutDataSet(&dataset)
	q.limit = limit
	if err != nil {
		return nil, nil, err
	}

	if dataset.Len() == 0 {
		return nil, dataset.Fields, err
	}
	if isPool {
		return dataset.RowAtP(0), dataset.Fields, nil
	}
	return dataset.RowAt(0), dataset.Fields, nil
}

func setQueryFields(q *QueryBuilder, buf *bytebufferpool.ByteBuffer) {
	fields := BStarKey
	if q.fields != "" {
		fields = []byte(q.fields)
	}

	buf.Write(fields)
}
