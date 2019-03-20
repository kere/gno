package db

import (
	"bytes"
	"fmt"
	"strconv"
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
	isExec     bool
	isQueryOne bool
}

// NewQuery new
func NewQuery(t string) *QueryBuilder {
	return &QueryBuilder{table: t, isExec: true}
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

// SetIsPrepare prepare sql
func (q *QueryBuilder) SetIsPrepare(v bool) *QueryBuilder {
	q.isExec = !v
	return q
}

// GetIsPrepare get
func (q *QueryBuilder) GetIsPrepare() bool {
	return !q.isExec
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

func querybuildCacheKey(q *QueryBuilder) string {
	s := bytes.Buffer{}
	s.WriteString(q.table)
	setQueryFields(q, &s)
	s.WriteString(fmt.Sprintf("%d%d", q.limit, q.offset))
	// s.WriteString(strconv.FormatInt(q.limit, 10) + strconv.Formart(q.offset, 10))
	s.WriteString(q.order)
	if q.where != "" {
		s.WriteString(q.where)
		s.WriteString(fmt.Sprint(q.args))
	}

	return fmt.Sprintf("db-%s:%x", q.GetDatabase().Name, MD5(s.Bytes()))
}

// ClearCache func
func (q *QueryBuilder) ClearCache() error {
	if cacheIns == nil {
		return nil
	}
	return cacheDel(querybuildCacheKey(q))
}

func parseQuery(q *QueryBuilder) string {
	s := bytes.Buffer{}
	s.Write(bSQLSelect)

	setQueryFields(q, &s)

	s.Write(bSQLFrom)
	s.WriteString(q.table)

	if len(q.leftJoin) > 0 {
		s.WriteString(" left join ")
		s.WriteString(q.leftJoin)
	}

	if q.where != "" {
		s.Write(bSQLWhere)
		s.WriteString(q.GetDatabase().Driver.Adapt(q.where, 0))
	}
	if q.order != "" {
		s.Write(bSQLOrder)
		s.WriteString(q.order)
	}

	limit := q.limit
	if q.isQueryOne {
		limit = 1
	}

	if limit > 0 && q.offset > 0 {
		s.Write(bSQLLimit)
		s.WriteString(string(limit))
		s.Write(bSQLOffset)
		s.WriteString(string(q.offset))
	} else if limit > 0 {
		s.Write(bSQLLimit)
		s.WriteString(strconv.FormatInt(int64(limit), 10))
	} else if q.offset > 0 {
		s.Write(bSQLOffset)
		s.WriteString(strconv.FormatInt(int64(q.offset), 10))
	}
	return s.String()
}

// Query return DataSet
func (q *QueryBuilder) Query() (DataSet, error) {
	var key string

	if q.cache {
		key = querybuildCacheKey(q)
		if exi, _ := cacheIns.IsExists(key); exi {
			return cacheGet(key)
		}
	}
	var r DataSet
	var err error
	if q.isExec {
		r, err = q.GetDatabase().Query(parseQuery(q), q.args...)
	} else {
		r, err = q.GetDatabase().QueryPrepare(parseQuery(q), q.args...)
	}

	if err != nil {
		return r, err
	}
	if q.cache {
		// r.Bytes2String()
		cacheSet(key, r, q.expire)
	}

	return r, nil
}

// QueryOne limit=1
func (q *QueryBuilder) QueryOne() (DataRow, error) {
	r, err := q.Query()
	if err != nil {
		return nil, err
	}

	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

// // Find return VODataSet
// func (q *QueryBuilder) Find() (VODataSet, error) {
// 	database := q.getDatabase()
//
// 	var key string
// 	if q.cache {
// 		key = q.cachekey()
// 		if exi, _ := cacheIns.Exists(key); exi {
// 			return cacheGetX(key, q.cls)
// 		}
// 	}
// 	var r VODataSet
// 	var err error
// 	if q.isExec {
// 		r, err = database.FindPrepare(q.cls, NewSqlState(parseQuery(q)))
// 	} else {
// 		r, err = database.Find(q.cls, NewSqlState(parseQuery(q)))
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	if q.cache {
// 		cacheSet(key, r, q.expire)
// 	}
// 	return r, nil
// }

// func (q *QueryBuilder) FindOne() (IVO, error) {
// 	q.isQueryOne = true
// 	r, err := q.Find()
// 	q.isQueryOne = false
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if len(r) > 0 {
// 		return r[0], nil
// 	}
// 	return nil, nil
// }

func setQueryFields(q *QueryBuilder, s *bytes.Buffer) {
	field := BStarKey
	if q.field != "" {
		field = []byte(q.field)
	} else if q.cls != nil {
		sm := NewStructConvert(q.cls)
		field = bytes.Join(sm.DBFields(), BCommaSplit)
	}

	s.Write(field)
}

// TxQuery transction
func (q *QueryBuilder) TxQuery(tx *Tx) (DataSet, error) {
	return tx.Query(parseQuery(q), q.args...)
}

// TxQueryOne transction
func (q *QueryBuilder) TxQueryOne(tx *Tx) (DataRow, error) {
	return tx.QueryOne(parseQuery(q), q.args...)
}
