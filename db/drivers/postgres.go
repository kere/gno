package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/valyala/bytebufferpool"
)

var (
	brHSTORE = []byte("\"=>\"")
	brJSON   = []byte("\":\"")
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

// Adapt f
func (p *Postgres) Adapt(sql string, n int) string {
	arr := strings.Split(sql, sQuestionMark)
	l := len(arr)
	if l == 0 {
		return sql
	}
	var s strings.Builder
	for i := 0; i < l-1; i++ {
		if arr[i] == "" {
			continue
		}
		s.WriteString(arr[i])
		s.Write(bDollar)
		s.WriteString(fmt.Sprint(i + 1 + n))
	}

	return s.String()
}

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
	return fmt.Sprint("select currval(pg_get_serial_sequence('", table, "','", pkey, "')) as count")
}

var (
	bytePool   bytebufferpool.Pool
	emptyArray = "{}"
)

func sliceToStore(v interface{}) string {
	val := reflect.ValueOf(v)
	l := val.Len()
	if l == 0 {
		return emptyArray
	}
	arr := make([]string, l)

	for i := 0; i < l; i++ {
		v := val.Index(i)
		arr[i] = fmt.Sprint(v.Interface())
	}
	return fmt.Sprint("{", strings.Join(arr, ","), "}")
}

// func Hstore()string{
// 		// case map[string]string:
// 			var valStr string
// 			hdata := v.(map[string]string)
// 			arr := make([]string, len(hdata))
// 			i := 0
// 			for kk, vv := range hdata {
// 				valStr = strings.Replace(fmt.Sprint(vv), "\"", "\\\"", -1)
// 				arr[i] = fmt.Sprint("\"", kk, "\"", "=>", "\"", valStr, "\"")
// 				i++
// 			}
// 			return fmt.Sprint(strings.Join(arr, ","))
// }

// StoreData for value
func (p *Postgres) StoreData(key string, v interface{}) interface{} {
	if v == nil {
		return nil
	}

	if len(key) > 5 && key[len(key)-5:] == subfixJSON {
		b, _ := json.Marshal(v)
		return b
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

// Strings
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

// Float64s arr
func (p *Postgres) Float64s(src []byte) ([]float64, error) {
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

// func (p *Postgres) ParseStringSlice(src []byte, ptr interface{}) error {
// 	src = bytes.Replace(src, b_BRACE_LEFT, b_BRACKET_LEFT, -1)
// 	src = bytes.Replace(src, b_BRACE_RIGHT, b_BRACKET_RIGHT, -1)
// 	src = bytes.Replace(src, b_Quote, b_DoubleQuote, -1)
//
// 	if err := json.Unmarshal(src, ptr); err != nil {
// 		return fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
// 	}
//
// 	return nil
// }

// // HStore db
// func (p *Postgres) HStore(src []byte) (map[string]string, error) {
// 	src = bytes.Replace(src, brHSTORE, brJSON, -1)
// 	src = append(b_BRACE_LEFT, src...)
// 	v := make(map[string]string)
//
// 	if err := json.Unmarshal(append(src, b_BRACE_RIGHT...), &v); err != nil {
// 		return nil, fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
// 	}
// 	return v, nil
// }

// ParseNumberSlice db number slice
func (p *Postgres) ParseNumberSlice(src []byte, ptr interface{}) error {
	if len(src) == 0 {
		return nil
	}

	src = bytes.Replace(src, b_BRACE_LEFT, b_BRACKET_LEFT, -1)
	src = bytes.Replace(src, b_BRACE_RIGHT, b_BRACKET_RIGHT, -1)
	src = bytes.Replace(src, bNaN, bZero, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return err
	}

	return nil
}

// ParseStringSlice db number slice
func (p *Postgres) ParseStringSlice(src []byte, ptr interface{}) error {
	src = bytes.Replace(src, b_BRACE_LEFT, b_BRACKET_LEFT, -1)
	src = bytes.Replace(src, b_BRACE_RIGHT, b_BRACKET_RIGHT, -1)
	src = bytes.Replace(src, b_Quote, b_DoubleQuote, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}

	return nil
}

// QuoteField f
func (p *Postgres) QuoteField(str string) string {
	return `"` + str + `"`
}

// QuoteFieldB f
func (p *Postgres) QuoteFieldB(s string) []byte {
	l := len(s)

	// l+8 : "name"=$1234
	arr := make([]byte, l+2, l+8)
	arr[0] = '"'
	for i := 0; i < l; i++ {
		arr[i+1] = s[i]
	}
	arr[l+1] = '"'

	return arr
}
