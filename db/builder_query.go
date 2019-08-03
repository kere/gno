package db

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/valyala/bytebufferpool"
)

var (
	// DefaultQueryCacheExpire redis cache
	DefaultQueryCacheExpire = 300
)

// QueryBuilder class
type QueryBuilder struct {
	builder

	table      string
	field      string
	leftJoin   string
	where      string
	args       []interface{}
	order      string
	limit      int
	offset     int
	cache      bool
	expire     int
	cls        IVO
	isPrepare  bool
	isQueryOne bool
}

// NewQuery new
func NewQuery(t string) *QueryBuilder {
	return &QueryBuilder{table: t, isPrepare: true}
}

// Table return string
func (q *QueryBuilder) Table(t string) *QueryBuilder {
	q.table = t
	return q
}

// GetTable return string
func (q *QueryBuilder) GetTable() string {
	return q.table
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
	q.field = s
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

// Cache property
func (q *QueryBuilder) Cache() *QueryBuilder {
	if cacheIns == nil {
		q.cache = false
		return q
	}

	q.expire = DefaultQueryCacheExpire
	q.cache = true
	return q
}

// DisableCache property
func (q *QueryBuilder) DisableCache() *QueryBuilder {
	q.cache = false
	return q
}

//CacheExpire Set query cache
// if expire value is 0, the cache expire will use DefaultQueryCacheExpire as cache expire.
// if expire valie low than 0, it will disable cache.
func (q *QueryBuilder) CacheExpire(expire int) *QueryBuilder {
	if cacheIns == nil || expire < 0 {
		q.cache = false
		return q
	}

	if expire == 0 {
		q.expire = DefaultQueryCacheExpire
	} else {
		q.expire = expire
	}

	q.cache = true
	return q
}

func querybuildCacheKey(q *QueryBuilder, datamode int) string {
	// s := bytes.Buffer{}
	buf := bytePool.Get()

	buf.WriteString(q.table)
	setQueryFields(q, buf)
	buf.WriteString(fmt.Sprintf("%d%d", q.limit, q.offset))
	// buf.WriteString(strconv.FormatInt(q.limit, 10) + strconv.Formart(q.offset, 10))
	buf.WriteString(q.order)
	if q.where != "" {
		buf.WriteString(q.where)
		buf.WriteString(fmt.Sprint(q.args))
	}
	src := buf.Bytes()
	bytePool.Put(buf)

	return fmt.Sprintf("db-%d-%s:%x", datamode, q.GetDatabase().Name, MD5(src))
}

// ClearCache func
func (q *QueryBuilder) ClearCache() {
	if cacheIns == nil {
		return
	}
	cacheDel(querybuildCacheKey(q, 1))
	cacheDel(querybuildCacheKey(q, 0))
}

func (q *QueryBuilder) Parse() (string, []interface{}) {
	// s := bytes.Buffer{}
	buf := bytePool.Get()
	buf.Write(bSQLSelect)

	setQueryFields(q, buf)

	buf.Write(bSQLFrom)
	buf.WriteString(q.table)

	if len(q.leftJoin) > 0 {
		buf.Write(bSQLLeftJoin)
		buf.WriteString(q.leftJoin)
	}

	if q.where != "" {
		buf.Write(bSQLWhere)
		buf.Write(q.GetDatabase().Driver.Adapt(q.where, 0))
	}
	if q.order != "" {
		buf.Write(bSQLOrder)
		buf.WriteString(q.order)
	}

	limit := q.limit
	if q.isQueryOne {
		limit = 1
	}

	if limit > 0 && q.offset > 0 {
		buf.Write(bSQLLimit)
		buf.WriteString(string(limit))
		buf.Write(bSQLOffset)
		buf.WriteString(string(q.offset))
	} else if limit > 0 {
		buf.Write(bSQLLimit)
		buf.WriteString(strconv.FormatInt(int64(limit), 10))
	} else if q.offset > 0 {
		buf.Write(bSQLOffset)
		buf.WriteString(strconv.FormatInt(int64(q.offset), 10))
	}
	str := buf.String()
	bytePool.Put(buf)
	return str, q.args
}

// Query return DataSet
func (q *QueryBuilder) Query() (DataSet, error) {
	dataset, _, err := q.cQuery(0)
	return dataset, err
}

// QueryRows return DataSet
func (q *QueryBuilder) QueryRows() (MapRows, error) {
	_, rows, err := q.cQuery(1)
	return rows, err
}

func (q *QueryBuilder) cQuery(mode int) (DataSet, MapRows, error) {
	var key string

	if q.cache {
		key = querybuildCacheKey(q, mode)
		if exi, _ := cacheIns.IsExists(key); exi {
			return cacheGet(key, mode)
		}
	}
	var dataset DataSet
	var datarows MapRows
	var err error

	sqlstr, _ := q.Parse()
	if mode == 1 {
		if q.isPrepare {
			datarows, err = q.GetDatabase().QueryRowsPrepare(sqlstr, q.args...)
		} else {
			datarows, err = q.GetDatabase().QueryRows(sqlstr, q.args...)
		}
	} else {
		if q.isPrepare {
			dataset, err = q.GetDatabase().QueryPrepare(sqlstr, q.args...)
		} else {
			dataset, err = q.GetDatabase().Query(sqlstr, q.args...)
		}
	}

	if err != nil {
		return dataset, datarows, err
	}

	if q.cache {
		if mode == 1 {
			cacheSet(key, datarows, q.expire)
		} else {
			cacheSet(key, dataset, q.expire)
		}
	}

	return dataset, datarows, nil
}

// QueryOne limit=1
func (q *QueryBuilder) QueryOne() (MapRow, error) {
	limit := q.limit
	q.limit = 1
	r, err := q.QueryRows()
	q.limit = limit
	if err != nil {
		return nil, err
	}

	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

func setQueryFields(q *QueryBuilder, buf *bytebufferpool.ByteBuffer) {
	field := BStarKey
	if q.field != "" {
		field = []byte(q.field)
	} else if q.cls != nil {
		sm := NewStructConvert(q.cls)
		field = bytes.Join(sm.DBFields(), BCommaSplit)
	}

	buf.Write(field)
}

// TxQuery transction
func (q *QueryBuilder) TxQuery(tx *Tx) (DataSet, error) {
	sqlstr, _ := q.Parse()
	if q.isPrepare {
		return tx.QueryPrepare(sqlstr, q.args...)
	}
	return tx.Query(sqlstr, q.args...)
}

// TxQueryRows transction
func (q *QueryBuilder) TxQueryRows(tx *Tx) (MapRows, error) {
	sqlstr, _ := q.Parse()
	if q.isPrepare {
		return tx.QueryRowsPrepare(sqlstr, q.args...)
	}
	return tx.QueryRows(sqlstr, q.args...)
}

// TxQueryOne transction
func (q *QueryBuilder) TxQueryOne(tx *Tx) (MapRow, error) {
	sqlstr, _ := q.Parse()
	if q.isPrepare {
		return tx.QueryOnePrepare(sqlstr, q.args...)
	}
	return tx.QueryOne(sqlstr, q.args...)
}
