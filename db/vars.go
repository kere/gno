package db

// 月份 1,01,Jan,January
// 日　 2,02,_2
// 时　 3,03,15,PM,pm,AM,am
// 分　 4,04
// 秒　 5,05
// 年　 06,2006
// 时区 -07,-0700,Z0700,Z07:00,-07:00,MST
// 周几 Mon,Monday

const (
	// ActionUpdate update
	ActionUpdate = 1
	// ActionInsert insert
	ActionInsert = 3
	// ActionDelete delete
	ActionDelete = 9

	// timeClassName s
	timeClassName = "time.Time"
	// ColumnBytePrefix prefix
	ColumnBytePrefix = "byte_"
	//FieldCount count
	FieldCount = "count"
	//FieldJSON json
	FieldJSON = "json"

	// SQuot ''
	SQuot = "'"

	// SLineBreak = "\n"
	SLineBreak = "\n"
	// SDoller $
	SDoller = "$"
)

var (
	bPGReturning = []byte(" RETURNING id")
	// sPGReturning = " RETURNING id"

	// DateTimeFormat = 2018-08-27 21:24:08.097823 +0000 GMT
	DateTimeFormat = "2006-01-02 15:04:05 -0700 MST"
	// DTFormat not with time zone
	DTFormat = "2006-01-02 15:04:05"

	bSQLSelect   = []byte("select ")
	bSQLDelete   = []byte("delete ")
	bSQLUpdate   = []byte("update ")
	bSQLSet      = []byte(" set ")
	bSQLFrom     = []byte(" from ")
	bSQLWhere    = []byte(" where ")
	bSQLOrder    = []byte(" order by ")
	bSQLLimit    = []byte(" limit ")
	bSQLLimitOne = []byte(" limit 1")
	bSQLOffset   = []byte(" offset ")

	// BDoller $
	BDoller = []byte("$")
	//BQuestionMark ?
	BQuestionMark = []byte("?")
	// BNull null
	BNull = []byte("NULL")
	// BEqual =
	BEqual = []byte("=")

	//BCommaSplit ,
	BCommaSplit = []byte(",")
	//SCommaSplit ,
	SCommaSplit = ","
	// BEmptyString ''
	BEmptyString = []byte("")
	//BStarKey *
	BStarKey = []byte("*")
	//BSpace ' '
	BSpace = []byte(" ")
)
