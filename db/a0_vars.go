package db

const (
	// SQuot ''
	SQuot = "'"

	//SComma ,
	SComma      = ","
	SBraceLeft  = "{"
	SBraceRight = "{"
	// SLineBreak = "\n"
	SLineBreak = "\n"
	// SDoller $
	SDoller = "$"

	// DateTimeFormat = 2018-08-27 21:24:08.097823 +0000 GMT
	DateTimeFormat = "2006-01-02 15:04:05 -0700 MST"
	// DTFormat not with time zone
	DTFormat = "2006-01-02 15:04:05"
)

var (
	bPGReturning = []byte(" RETURNING id")
	// sPGReturning = " RETURNING id"

	bSQLSelect   = []byte("SELECT ")
	bSQLDelete   = []byte("DELETE ")
	bSQLUpdate   = []byte("UPDATE ")
	bSQLSet      = []byte(" SET ")
	bSQLFrom     = []byte(" FROM ")
	bSQLWhere    = []byte(" WHERE ")
	bSQLOrder    = []byte(" ORDER BY ")
	bSQLLimit    = []byte(" LIMIT ")
	bSQLLimitOne = []byte(" LIMIT 1")
	bSQLOffset   = []byte(" OFFSET ")
	bSQLLeftJoin = []byte(" as a LEFT JOIN ")

	bInsSQL      = []byte("INSERT INTO ")
	bInsBracketL = []byte(" (")
	bInsBracketR = []byte(") VALUES ")

	// BDoller $
	BDoller = []byte("$")
	// //BQuestionMark ?
	// BQuestionMark = []byte("?")
	// BNull null
	BNull = []byte("NULL")
	// BEqual =
	BEqual = []byte("=")

	//BComma ,
	BComma = []byte(",")

	// BEmptyString ''
	BEmptyString = []byte("")
	//BStarKey *
	BStarKey = []byte("*")
	//BSpace ' '
	BSpace        = []byte(" ")
	BBraceLeft    = []byte("{")
	BBraceRight   = []byte("}")
	BBRACKET_LEFT = []byte("[")
	BBracketLeft  = []byte("[")
	BBracketRight = []byte("]")
	// bDollar         = []byte("$")
	BQuestionMark = []byte("?")
	BQuote        = []byte("'")
	BDoubleQuote  = []byte("\"")
	bZero         = []byte("0")
	BNaN          = []byte("NaN")
)
