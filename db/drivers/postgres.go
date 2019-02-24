package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	b_r_HSTORE = []byte("\"=>\"")
	b_r_JSON   = []byte("\":\"")
)

//Postgres class
type Postgres struct {
	Common
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

	} else {
		return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
			p.DBName,
			p.User,
			p.Password,
			p.Host,
			p.Port)
	}
}

func (p *Postgres) QuoteField(s string) string {
	return fmt.Sprint("\"", s, "\"")
}

func (p *Postgres) LastInsertId(table, pkey string) string {
	// return "select currval(pg_get_serial_sequence('" + table + "','" + pkey + "'))"
	return fmt.Sprint("select currval(pg_get_serial_sequence('", table, "','", pkey, "')) as count")
}

func (p *Postgres) sliceToStore(typ reflect.Type, v interface{}) string {
	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		value := reflect.ValueOf(v)
		arr := make([]string, value.Len())
		l := value.Len()
		if l == 0 {
			return "{}"
		}

		var tmpV reflect.Value
		for i := 0; i < l; i++ {
			tmpV = value.Index(i)
			arr[i] = p.sliceToStore(tmpV.Type(), tmpV.Interface())
		}
		return fmt.Sprint("{", strings.Join(arr, ","), "}")

	case reflect.String:
		return fmt.Sprint("'", v, "'")

	default:
		return fmt.Sprint(v)

	}

}

// FlatData for value
func (p *Postgres) FlatData(typ reflect.Type, v interface{}) interface{} {
	if v == nil {
		return "NULL"
	}

	switch typ.Kind() {
	// case reflect.String, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
	default:
		return v

	case reflect.Bool:
		if v.(bool) {
			return "t"
		} else {
			return "f"
		}

	case reflect.Array:
		return p.sliceToStore(typ, v)

	case reflect.Slice:
		switch v.(type) {
		case []byte:
			return v
		default:
			return p.sliceToStore(typ, v)
		}

	case reflect.Map:
		switch v.(type) {
		case map[string]string:
			var valStr string
			hdata := v.(map[string]string)
			arr := make([]string, len(hdata))
			i := 0
			for kk, vv := range hdata {
				valStr = strings.Replace(fmt.Sprint(vv), "\"", "\\\"", -1)
				arr[i] = fmt.Sprint("\"", kk, "\"", "=>", "\"", valStr, "\"")
				i++
			}
			return fmt.Sprint(strings.Join(arr, ","))

		default:
			b, err := json.Marshal(v)
			if err != nil {
				return []byte("")
			}
			return b

		}

	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			return v

		default:
			b, err := json.Marshal(v)
			if err != nil {
				return []byte("")
			}
			return b
		}
	}

}

func (p *Postgres) StringSlice(src []byte) ([]string, error) {
	if len(src) == 0 {
		return []string{}, nil
	}

	src = bytes.TrimPrefix(src, b_BRACE_LEFT)
	src = bytes.TrimSuffix(src, b_BRACE_RIGHT)
	if len(src) == 0 {
		return []string{}, nil
	}

	arr := bytes.Split(src, b_COMMA)
	l := len(arr)
	v := make([]string, len(arr))
	for i := 0; i < l; i++ {
		v[i] = string(bytes.Trim(arr[i], "'"))
	}

	return v, nil
}

func (p *Postgres) Int64Slice(src []byte) ([]int64, error) {
	if len(src) == 0 {
		return []int64{}, nil
	}
	var arr = make([]int64, 0)
	if err := p.ParseNumberSlice(src, &arr); err != nil {
		return nil, err
	}

	return arr, nil
}

func (p *Postgres) ParseStringSlice(src []byte, ptr interface{}) error {
	src = bytes.Replace(src, b_BRACE_LEFT, b_BRACKET_LEFT, -1)
	src = bytes.Replace(src, b_BRACE_RIGHT, b_BRACKET_RIGHT, -1)
	src = bytes.Replace(src, b_Quote, b_DoubleQuote, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}

	return nil
}

// HStore db
func (p *Postgres) HStore(src []byte) (map[string]string, error) {
	src = bytes.Replace(src, b_r_HSTORE, b_r_JSON, -1)
	src = append(b_BRACE_LEFT, src...)
	v := make(map[string]string)

	if err := json.Unmarshal(append(src, b_BRACE_RIGHT...), &v); err != nil {
		return nil, fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}
	return v, nil
}

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
