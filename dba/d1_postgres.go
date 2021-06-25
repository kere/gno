package dba

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/kere/gno/libs/util"
	"github.com/lib/pq"
	"github.com/valyala/bytebufferpool"
)

const (
	emptyArray = "{}"
)

var (
	bytePool bytebufferpool.Pool
)

//Postgres class
type Postgres struct {
	DBName   string
	User     string
	Password string
	Host     string
	HostAddr string
	Port     string
}

// Name f
func (p *Postgres) Name() string {
	return DriverPSQL
}

// // Adapt f
// func (p *Postgres) Adapt(sqlstr string, n int) string {
// 	src := util.Str2Bytes(sqlstr)
// 	arr := bytes.Split(src, BQuestionMark)
// 	// arr := strings.Split(sqlstr, sQuestionMark)
// 	l := len(arr)
// 	if l == 0 {
// 		return ""
// 	}
// 	if l == 1 {
// 		return sqlstr
// 	}
//
// 	buf := bytePool.Get()
// 	for i := 0; i < l-1; i++ {
// 		buf.Write(arr[i])
// 		buf.WriteByte('$')
//
// 		buf.WriteString(fmt.Sprint(i + 1 + n))
// 	}
//
// 	b := buf.String()
// 	bytePool.Put(buf)
//
// 	return b
// }

// ConnectString f
func (p *Postgres) ConnectString() string {
	if p.Host == "" {
		p.Host = "127.0.0.1"
	}
	if p.Port == "" {
		p.Port = "5432"
	}

	if p.HostAddr != "" {
		return fmt.Sprintf("dbname=%s user=%s password=%s hostaddr=%s sslmode=disable",
			p.DBName,
			p.User,
			p.Password,
			p.HostAddr)

	}
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		p.DBName,
		p.User,
		p.Password,
		p.Host,
		p.Port)
}

// LastInsertID f
func (p *Postgres) LastInsertID(table, pkey string) string {
	// return "select currval(pg_get_serial_sequence('" + table + "','" + pkey + "'))"
	return fmt.Sprint("SELECT currval(pg_get_serial_sequence('", table, "','", pkey, "')) as count")
}

func sliceToStore(v interface{}) string {
	val := reflect.ValueOf(v)
	l := val.Len()
	if l == 0 {
		return emptyArray
	}
	arr := make([]string, 0, l)

	for i := 0; i < l; i++ {
		v := fmt.Sprint(val.Index(i).Interface())
		if v == "" {
			continue
		}
		arr = append(arr, v)
	}
	return fmt.Sprint("{", strings.Join(arr, ","), "}")
}

// StoreData for value
func (p *Postgres) StoreData(key string, v interface{}) interface{} {
	if v == nil {
		return nil
	}

	typ := reflect.TypeOf(v)
	switch typ.Kind() {
	default:
		return v

	case reflect.Array:
		return sliceToStore(v)

	case reflect.Slice:
		switch v.(type) {
		case []byte:
			return v
		default:
			return sliceToStore(v)
		}

	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			return v

		default:
			b, _ := json.Marshal(v)
			return b
		}
	}
}

// Strings []string
func (p *Postgres) Strings(src []byte) ([]string, error) {
	if len(src) == 0 {
		return nil, nil
	}

	arr := pq.StringArray{}
	return arr, arr.Scan(src)
}

// Int64s arr
func (p *Postgres) Int64s(src []byte) ([]int64, error) {
	if len(src) == 0 {
		return nil, nil
	}

	arr := pq.Int64Array{}
	return arr, arr.Scan(src)
}

// Floats arr
func (p *Postgres) Floats(src []byte) ([]float64, error) {
	if len(src) == 0 {
		return nil, nil
	}

	arr := pq.Float64Array{}
	return arr, arr.Scan(src)
}

// Ints arr
func (p *Postgres) Ints(src []byte) ([]int, error) {
	if len(src) == 0 {
		return nil, nil
	}
	var vals []int
	err := p.ParseNumberSlice(src, &vals)
	return vals, err
}

// ParseNumberSlice db number slice
func (p *Postgres) ParseNumberSlice(src []byte, ptr interface{}) error {
	if len(src) == 0 {
		return nil
	}

	src = bytes.Replace(src, BBRACE_LEFT, BBRACKET_LEFT, -1)
	src = bytes.Replace(src, BBRACE_RIGHT, BBRACKET_RIGHT, -1)
	src = bytes.Replace(src, BNaN, bZero, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return err
	}

	return nil
}

// ParseStringSlice db number slice
func (p *Postgres) ParseStringSlice(src []byte, ptr interface{}) error {
	src = bytes.Replace(src, BBRACE_LEFT, BBRACKET_LEFT, -1)
	src = bytes.Replace(src, BBRACE_RIGHT, BBRACKET_RIGHT, -1)
	src = bytes.Replace(src, BQuote, BDoubleQuote, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}

	return nil
}

// WriteQuoteIdentifier f
func (p *Postgres) WriteQuoteIdentifier(w io.Writer, s string) {
	// return pq.QuoteLiteral(literal)
	// `"` + strings.Replace(name, `"`, `""`, -1) + `"`
	str := strings.Replace(s, `"`, `""`, -1)
	w.Write(BDoubleQuote)
	w.Write(util.Str2Bytes(str))
	w.Write(BDoubleQuote)
}

// // QuoteLiteral quotes a 'literal' (e.g. a parameter, often used to pass literal
// // to DDL and other statements that do not accept parameters) to be used as part
// // of an SQL statement.  For example:
// //
// //    exp_date := pq.QuoteLiteral("2023-01-05 15:00:00Z")
// //    err := db.Exec(fmt.Sprintf("CREATE ROLE my_user VALID UNTIL %s", exp_date))
// //
// // Any single quotes in name will be escaped. Any backslashes (i.e. "\") will be
// // replaced by two backslashes (i.e. "\\") and the C-style escape identifier
// // that PostgreSQL provides ('E') will be prepended to the string.
// func (p *Postgres) QuoteLiteral(literal string) string {
// 	return pq.QuoteLiteral(literal)
// }
