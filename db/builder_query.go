package db

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var (
	// DefaultQueryCacheExpire redis cache
	DefaultQueryCacheExpire = 300
)

// Querybuilder class
type QueryBuilder struct {
	builder

	table      string
	field      string
	leftJoin   string
	unselect   []string
	where      *CondParams
	order      string
	limit      int
	offset     int
	cache      bool
	expire     int
	cls        IVO
	isPrepare  bool
	isQueryOne bool
}

// NewQueryBuilder new
func NewQueryBuilder(t string) *QueryBuilder {
	return (&QueryBuilder{}).Table(t).IsPrepare(true)
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

func (q *QueryBuilder) IsPrepare(v bool) *QueryBuilder {
	q.isPrepare = v
	return q
}

func (q *QueryBuilder) GetIsPrepare() bool {
	return q.isPrepare
}

func (q *QueryBuilder) Select(s string) *QueryBuilder {
	q.field = s
	return q
}

func (q *QueryBuilder) LeftJoin(s string) *QueryBuilder {
	q.leftJoin = s
	return q
}

func (q *QueryBuilder) UnSelect(s ...string) *QueryBuilder {
	q.unselect = s
	return q
}

func (q *QueryBuilder) Where(s string, args ...interface{}) *QueryBuilder {
	if s == "" {
		return q
	}

	q.where = &CondParams{s, args}
	return q
}

func (q *QueryBuilder) Order(s string) *QueryBuilder {
	q.order = s
	return q
}

func (q *QueryBuilder) Page(page int, pageSize int) *QueryBuilder {
	q.offset = (page - 1) * pageSize
	q.limit = pageSize
	return q
}

func (q *QueryBuilder) Limit(n int) *QueryBuilder {
	q.limit = n
	return q
}

func (q *QueryBuilder) Offset(n int) *QueryBuilder {
	q.offset = n
	return q
}

func (q *QueryBuilder) Struct(cls IVO) *QueryBuilder {
	q.cls = cls
	return q
}

func (q *QueryBuilder) Cache() *QueryBuilder {
	if cacheIns == nil {
		q.cache = false
		return q
	}

	q.expire = DefaultQueryCacheExpire
	q.cache = true
	return q
}

func (q *QueryBuilder) DisableCache() *QueryBuilder {
	q.cache = false
	return q
}

// Set query cache
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

func (q *QueryBuilder) cachekey() string {
	s := bytes.Buffer{}
	s.WriteString(q.table)
	q.writeField(&s)
	s.WriteString(fmt.Sprintf("%d%d", q.limit, q.offset))
	// s.WriteString(strconv.FormatInt(q.limit, 10) + strconv.Formart(q.offset, 10))

	s.WriteString(q.order)
	if q.where != nil {
		s.WriteString(fmt.Sprintf("%s%v", q.where.Cond, q.where.Args))
	}
	return string(MD5(s.Bytes()))
}

func (q *QueryBuilder) ClearCache() error {
	if cacheIns == nil {
		return nil
	}
	return cacheDel(q.cachekey())
}

func (q *QueryBuilder) parse() ([]byte, []interface{}) {
	s := bytes.Buffer{}
	var args []interface{}
	if q.getDatabase() == nil {
		panic("database is nil, check db Configuration!")
	}
	driver := q.getDatabase().Driver

	s.Write(bSQLSelect)

	q.writeField(&s)

	s.Write(bSQLFrom)
	s.WriteString(driver.QuoteField(q.table))

	if len(q.leftJoin) > 0 {
		s.WriteString(" left join ")
		s.WriteString(q.leftJoin)
	}

	if q.where != nil {
		s.Write(bSQLWhere)
		s.WriteString(q.where.Cond)
		args = q.where.Args
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
	return s.Bytes(), args
}

// Query return DataSet
func (q *QueryBuilder) Query() (DataSet, error) {
	database := q.getDatabase()

	var key string
	if q.cache {
		key = q.cachekey()
		if exi, _ := cacheIns.Exists(key); exi {
			return cacheGet(key)
		}
	}
	var r DataSet
	var err error
	if q.isPrepare {
		r, err = database.QueryPrepare(NewSqlState(q.parse()))
	} else {
		r, err = database.Query(NewSqlState(q.parse()))
	}

	if err != nil {
		return r, err
	}
	if q.cache {
		r.Bytes2String()
		cacheSet(key, r, q.expire)
	}

	return r, nil
}

// QueryOne limit=1
func (q *QueryBuilder) QueryOne() (DataRow, error) {
	q.isQueryOne = true
	r, err := q.Query()
	q.isQueryOne = false
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
// 	if q.isPrepare {
// 		r, err = database.FindPrepare(q.cls, NewSqlState(q.parse()))
// 	} else {
// 		r, err = database.Find(q.cls, NewSqlState(q.parse()))
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

func (q *QueryBuilder) writeField(s *bytes.Buffer) {
	field := B_StarKey
	if q.field != "" {
		if len(q.unselect) == 0 {
			field = []byte(q.field)

		} else {
			arr := strings.Split(q.field, ",")
			tmp := make([]string, 0)
			for _, item := range arr {
				if InStrings(q.unselect, item) {
					continue
				}
				tmp = append(tmp, item)
			}
			field = []byte(strings.Join(tmp, ","))
		}

	} else if q.cls != nil {
		sm := NewStructConvert(q.cls)
		sm.SetExcludes(q.unselect)
		field = bytes.Join(sm.DBFields(), B_CommaSplit)

	}

	s.Write(field)
}

func (q *QueryBuilder) TxQuery(tx *Tx) (DataSet, error) {
	return tx.Query(NewSqlState(q.parse()))
}

func (q *QueryBuilder) TxQueryOne(tx *Tx) (DataRow, error) {
	q.limit = 1
	return tx.QueryOne(NewSqlState(q.parse()))
}
