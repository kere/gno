package drivers

import "reflect"

const (
	// DriverPSQL pgsql
	DriverPSQL = "postgres"

	// DriverMySQL mysql
	DriverMySQL = "mysql"
	// DriverSqlite sqlite
	DriverSqlite = "sqlite"

	sQuestionMark = "?"
)

var (
	b_BRACE_LEFT    = []byte("{")
	b_BRACE_RIGHT   = []byte("}")
	b_BRACKET_LEFT  = []byte("[")
	b_BRACKET_RIGHT = []byte("]")
	b_COMMA         = []byte(",")
	bDollar         = []byte("$")
	BQuestionMark   = []byte("?")
	b_Quote         = []byte("'")
	b_DoubleQuote   = []byte("\"")
	bZero           = []byte("0")
	bNaN            = []byte("NaN")
)

// Common class
type Common struct {
	connect string
}

func (c *Common) Adapt(b string, n int) string {
	return b
}

func (c *Common) SetConnectString(s string) {
	c.connect = s
}

func (c *Common) ConnectString() string {
	return c.connect
}

func (c *Common) QuoteField(s string) string {
	return `"` + s + `"`
}

func (c *Common) QuoteFieldB(s string) []byte {
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

func (c *Common) LastInsertID(table, pkey string) string {
	return ""
}

func (c *Common) Int64Slice(str []byte) ([]int64, error) {
	return nil, nil
}

func (c *Common) ParseStringSlice(src []byte, ptr interface{}) error {
	return nil
}

func (c *Common) FlatData(typ reflect.Type, v interface{}) interface{} {
	return nil
}

func (c *Common) ParseNumberSlice(src []byte, ptr interface{}) error {
	return nil
}

func (c *Common) StringSlice(src []byte) ([]string, error) {
	return nil, nil
}

func (c *Common) HStore(src []byte) (map[string]string, error) {
	return nil, nil
}

func (c *Common) Name() string {
	return ""
}
