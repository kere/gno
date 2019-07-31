package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kere/gno/db/drivers"
	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/myerr"
)

var (
	dbpool *databasePool
	// dbConf conf.Conf
)

func init() {
	dbpool = &databasePool{dblist: make(map[string]*Database)}
}

type databasePool struct {
	dblist  map[string]*Database
	current *Database
}

// Current database
func (dp *databasePool) Current() *Database {
	return dp.current
}

// SetCurrent database
func (dp *databasePool) SetCurrent(d *Database) {
	dp.current = d
}

// Use database
func (dp *databasePool) Use(name string) {
	c := dp.GetDatabase(name)
	if c == nil {
		fmt.Println(name, " database is not found!")
		return
	}
	dp.current = c
}

// SetDatabase by name
func (dp *databasePool) SetDatabase(name string, d *Database) {
	dp.dblist[name] = d
}

// GetDatabase by name
func (dp *databasePool) GetDatabase(name string) *Database {
	if v, ok := dp.dblist[name]; ok {
		return v
	}
	return nil
}

// Init it
func Init(name string, config map[string]string) {
	fmt.Println("Init Database", config)
	dbpool.SetCurrent(New(name, config))
}

func confGet(config map[string]string, key string) string {
	if v, ok := config[key]; ok {
		return v
	}
	return ""
}

// New func
// create a database instance
func New(name string, c map[string]string) *Database {
	if dbpool.GetDatabase(name) != nil {
		panic(name + " this database is already exists!")
	}

	if c == nil {
		return nil
	}

	driverName := confGet(c, "driver")
	logger := NewLogger(c)

	var driver IDriver
	switch driverName {
	case "postgres", "psql":
		driver = &drivers.Postgres{DBName: confGet(c, "dbname"),
			User:     confGet(c, "user"),
			Password: confGet(c, "password"),
			Host:     confGet(c, "host"),
			HostAddr: confGet(c, "hostaddr"),
			Port:     confGet(c, "port"),
		}

	case "mysql":
		driver = &drivers.Mysql{DBName: confGet(c, "dbname"),
			User:       confGet(c, "user"),
			Password:   confGet(c, "password"),
			Protocol:   confGet(c, "protocol"),
			Parameters: confGet(c, "parameters"),
			Addr:       confGet(c, "addr")}

	case "sqlite3":
		driver = &drivers.Sqlite3{File: confGet(c, "file")}

	default:
		logger.Println("you may need regist a custom driver: db.RegistDriver(Mysql{})")
		driver = &drivers.Common{}

	}

	driver.SetConnectString(confGet(c, "connect"))
	d := NewDatabase(name, driver, conf.Conf(c), logger)

	dbpool.SetDatabase(name, d)
	dbpool.SetCurrent(d)
	return d
}

func logSQLErr(sqlstr string, args []interface{}) {
	sep := ": "
	var s strings.Builder
	s.WriteString(sqlstr)
	s.WriteString(SLineBreak)
	l := len(args)
	for i := 0; i < l; i++ {
		s.WriteString(fmt.Sprint(i, sep))

		switch args[i].(type) {
		case []byte:
			s.Write(args[i].([]byte))
		default:
			s.WriteString(fmt.Sprint(args[i]))
		}
		s.WriteString(SLineBreak)
	}
	s.WriteString(SLineBreak)
	log.App.Error(s.String())
}

// CQuery from database on prepare mode
func CQuery(conn *sql.DB, sqlstr string, args ...interface{}) (DataSet, error) {
	dataset, _, err := cQuery(0, 0, conn, sqlstr, args...)
	return dataset, err
}

// CQueryPrepare from database on prepare mode
func CQueryPrepare(conn *sql.DB, sqlstr string, args ...interface{}) (DataSet, error) {
	result, _, err := cQuery(0, 1, conn, sqlstr, args...)
	return result, err
}

// CQueryRows from database on prepare mode
func CQueryRows(conn *sql.DB, sqlstr string, args ...interface{}) (MapRows, error) {
	_, rows, err := cQuery(1, 0, conn, sqlstr, args...)
	return rows, err
}

// CQueryPrepareRows from database on prepare mode
func CQueryPrepareRows(conn *sql.DB, sqlstr string, args ...interface{}) (MapRows, error) {
	_, rows, err := cQuery(1, 1, conn, sqlstr, args...)
	return rows, err
}

// pmode 1: prepare
// result mode 1: Rows
func cQuery(mode, pmode int, conn *sql.DB, sqlstr string, args ...interface{}) (DataSet, MapRows, error) {
	var rows *sql.Rows
	var err error
	var stmt *sql.Stmt
	var dataset DataSet

	defer conn.Close()

	if pmode == 1 {
		stmt, err = conn.Prepare(sqlstr)
		defer stmt.Close()
		if err != nil {
			logSQLErr(sqlstr, args)
			return dataset, nil, myerr.New(err).Log().Stack()
		}
		rows, err = stmt.Query(args...)
		defer rows.Close()
	} else {
		rows, err = conn.Query(sqlstr, args...)
		defer rows.Close()
	}

	if err != nil {
		logSQLErr(sqlstr, args)
		return dataset, nil, myerr.New(err).Log().Stack()
	}

	var maprows MapRows
	if mode == 1 {
		maprows, err = ScanToMapRows(rows)
	} else {
		dataset, err = ScanToDataSet(rows)
	}

	if err != nil {
		return dataset, nil, myerr.New(err).Log().Stack()
	}

	return dataset, maprows, nil
}

// Exec sql.
// If your has more than on sql command, it will only excute the first.
// This function use the current database from database bool
func Exec(conn *sql.DB, sqlstr string, args ...interface{}) (result sql.Result, err error) {
	defer conn.Close()
	if len(args) == 0 {
		result, err = conn.Exec(sqlstr)
	} else {
		result, err = conn.Exec(sqlstr, args...)
	}
	if err != nil {
		logSQLErr(sqlstr, args)
		return result, myerr.New(err).Log().Stack()
	}
	return result, nil
}

// ExecPrepare sql on prepare mode
// This function use the current database from database bool
func ExecPrepare(conn *sql.DB, sqlstr string, args ...interface{}) (result sql.Result, err error) {
	defer conn.Close()
	s, err := conn.Prepare(sqlstr)
	if err != nil {
		return nil, myerr.New(err).Log().Stack()
	}

	defer s.Close()

	if len(args) == 0 {
		result, err = s.Exec()
	} else {
		result, err = s.Exec(args...)
	}
	if err != nil {
		logSQLErr(sqlstr, args)
		return result, myerr.New(err).Log().Stack()
	}
	return result, nil
}

// Get a database instance by name from database pool
func Get(name string) *Database {
	return dbpool.GetDatabase(name)
}

//Current Return the current database from database pool
func Current() *Database {
	if dbpool.Current() == nil {
		panic("db is not initalized")
	}
	return dbpool.Current()
}

// Use current database by name
func Use(name string) {
	dbpool.Use(name)
}

// DatabaseCount Get database count
func DatabaseCount() int {
	return len(dbpool.dblist)
}
