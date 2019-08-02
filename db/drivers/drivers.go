package drivers

const (
	subfixJSON = "_json"

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
	// bDollar         = []byte("$")
	BQuestionMark = []byte("?")
	b_Quote       = []byte("'")
	b_DoubleQuote = []byte("\"")
	bZero         = []byte("0")
	bNaN          = []byte("NaN")
)
