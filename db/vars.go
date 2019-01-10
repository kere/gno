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
	ActionUpdate = "update"
	// ActionInsert insert
	ActionInsert = "insert"

	timeClassName = "time.Time"

	ColumnBytePrefix = "byte_"
)

var (
	// DateTimeFormat = 2018-08-27 21:24:08.097823 +0000 GMT
	DateTimeFormat = "2006-01-02 15:04:05 -0700 MST"

	bSQLSelect = []byte("select ")
	bSQLDelete = []byte("delete ")
	bSQLUpdate = []byte("update ")
	bSQLSet    = []byte(" set ")
	bSQLFrom   = []byte(" from ")
	bSQLWhere  = []byte(" where ")
	bSQLOrder  = []byte(" order by ")
	bSQLLimit  = []byte(" limit ")
	bSQLOffset = []byte(" offset ")

	B_QuestionMark = []byte("?")
	BNull          = []byte("NULL")
	B_Equal        = []byte("=")

	BCommaSplit   = []byte(",")
	B_EmptyString = []byte("")
	B_StarKey     = []byte("*")
	B_Space       = []byte(" ")
)
