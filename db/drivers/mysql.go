package drivers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Mysql struct {
	Common
	DBName     string
	User       string
	Password   string
	Addr       string
	Protocol   string
	Parameters string
}

// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
// A DSN in its fullest form:
// username:password@protocol(address)/dbname?param=value
// root:pw@unix(/tmp/mysql.sock)/myDatabase?loc=Local
// user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
func (this *Mysql) ConnectString() string {
	protocol := "tcp"
	if this.Protocol != "" {
		protocol = this.Protocol
	}

	addr := "127.0.0.1:3306"
	if this.Addr != "" {
		addr = this.Addr
	}

	if this.Parameters != "" {
		this.Parameters = "?" + this.Parameters
	}

	return fmt.Sprintf("%s:%s@%s(%s)/%s%s", this.User, this.Password, protocol, addr, this.DBName, this.Parameters)
}

func (this *Mysql) QuoteField(s string) string {
	return fmt.Sprint("`", s, "`")
}

func (this *Mysql) LastInsertId(table, id string) string {
	return "SELECT LAST_INSERT_ID() as count"
}

func (this *Mysql) DriverName() string {
	return DriverMySQL
}

func (this *Mysql) toJson(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("")
	}
	return b
}

func (this *Mysql) FlatData(typ reflect.Type, v interface{}) interface{} {
	if v == nil {
		return nil
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
		return this.toJson(v)

	case reflect.Slice:
		switch v.(type) {
		case []byte:
			return v
		default:
			return this.toJson(v)
		}

	case reflect.Map:
		return this.toJson(v)

	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			// return (v.(time.Time)).Format("2006-01-02 15:03:04")
			return v

		default:
			return this.toJson(v)
		}
	}

}

func (this *Mysql) StringSlice(src []byte) ([]string, error) {
	if len(src) == 0 {
		return []string{}, nil
	}

	src = bytes.TrimPrefix(src, b_BRACKET_LEFT)
	src = bytes.TrimSuffix(src, b_BRACKET_RIGHT)
	if len(src) == 0 {
		return []string{}, nil
	}

	l := bytes.Split(src, b_COMMA)
	v := make([]string, len(l))
	for i, _ := range l {
		v[i] = string(bytes.Trim(l[i], "'"))
	}

	return v, nil
}

func (this *Mysql) Int64Slice(src []byte) ([]int64, error) {
	if len(src) == 0 {
		return []int64{}, nil
	}
	var arr = make([]int64, 0)
	if err := this.ParseNumberSlice(src, &arr); err != nil {
		return nil, err
	}

	return arr, nil
}

func (this *Mysql) ParseStringSlice(src []byte, ptr interface{}) error {
	src = bytes.Replace(src, b_Quote, b_DoubleQuote, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}

	return nil
}

func (this *Mysql) ParseNumberSlice(src []byte, ptr interface{}) error {
	if src == nil {
		return errors.New("empty source")
	}

	// src = bytes.Replace(src, b_BRACE_LEFT, b_BRACKET_LEFT, -1)
	// src = bytes.Replace(src, b_BRACE_RIGHT, b_BRACKET_RIGHT, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}

	return nil
}

func (this *Mysql) HStore(src []byte) (map[string]string, error) {
	src = bytes.Replace(src, b_r_HSTORE, b_r_JSON, -1)
	src = append(b_BRACE_LEFT, src...)
	v := make(map[string]string)

	if err := json.Unmarshal(append(src, b_BRACE_RIGHT...), &v); err != nil {
		return nil, fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}
	return v, nil
}
