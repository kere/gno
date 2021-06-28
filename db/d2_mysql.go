package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/kere/gno/libs/util"
)

type Mysql struct {
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
func (m *Mysql) ConnectString() string {
	protocol := "tcp"
	if m.Protocol != "" {
		protocol = m.Protocol
	}

	addr := "127.0.0.1:3306"
	if m.Addr != "" {
		addr = m.Addr
	}

	if m.Parameters != "" {
		m.Parameters = "?" + m.Parameters
	}

	return fmt.Sprintf("%s:%s@%s(%s)/%s%s", m.User, m.Password, protocol, addr, m.DBName, m.Parameters)
}

func (m *Mysql) QuoteField(str string) string {
	return `"` + str + `"`
}

func (m *Mysql) QuoteFieldB(str string) []byte {
	l := len(str)

	// l+8 : "name"=$1234
	arr := make([]byte, l+2, l+8)
	arr[0] = '"'
	for i := 0; i < l; i++ {
		arr[i+1] = str[i]
	}
	arr[l+1] = '"'

	return arr
}

func (m *Mysql) LastInsertId(table, id string) string {
	return "SELECT LAST_INSERT_ID() as count"
}

func (m *Mysql) Name() string {
	return DriverMySQL
}

func (m *Mysql) toJson(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("")
	}
	return b
}

func (m *Mysql) FlatData(typ reflect.Type, v interface{}) interface{} {
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
		return m.toJson(v)

	case reflect.Slice:
		switch v.(type) {
		case []byte:
			return v
		default:
			return m.toJson(v)
		}

	case reflect.Map:
		return m.toJson(v)

	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			// return (v.(time.Time)).Format("2006-01-02 15:03:04")
			return v

		default:
			return m.toJson(v)
		}
	}
}

func (m *Mysql) StringSlice(src []byte) ([]string, error) {
	if len(src) == 0 {
		return []string{}, nil
	}

	src = bytes.TrimPrefix(src, util.BBracketLeft)
	src = bytes.TrimSuffix(src, util.BBracketRight)
	if len(src) == 0 {
		return []string{}, nil
	}

	l := bytes.Split(src, util.BComma)
	v := make([]string, len(l))
	for i, _ := range l {
		v[i] = string(bytes.Trim(l[i], "'"))
	}

	return v, nil
}

func (m *Mysql) Int64Slice(src []byte) ([]int64, error) {
	if len(src) == 0 {
		return []int64{}, nil
	}
	var arr = make([]int64, 0)
	if err := m.ParseNumberSlice(src, &arr); err != nil {
		return nil, err
	}

	return arr, nil
}

func (m *Mysql) ParseStringSlice(src []byte, ptr interface{}) error {
	src = bytes.Replace(src, util.BQuote, util.BDoubleQuote, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}

	return nil
}

func (m *Mysql) ParseNumberSlice(src []byte, ptr interface{}) error {
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

// func (m *Mysql) HStore(src []byte) (map[string]string, error) {
// 	src = bytes.Replace(src, brHSTORE, brJSON, -1)
// 	src = append(b_BRACE_LEFT, src...)
// 	v := make(map[string]string)
//
// 	if err := json.Unmarshal(append(src, b_BRACE_RIGHT...), &v); err != nil {
// 		return nil, fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
// 	}
// 	return v, nil
// }
