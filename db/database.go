package db

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/kere/gno/libs/log"
)

type SqlState struct {
	sql  []byte
	args []interface{}
}

func NewSqlState(sql []byte, args ...interface{}) *SqlState {
	ss := &SqlState{}
	ss.sql = Current().Driver.AdaptSql(sql)
	// for i, _ := range args {
	// 	args[i] = Current().Driver.FlatData(args[i])
	// }
	ss.args = args

	return ss
}

func (s *SqlState) SetSql(b []byte) {
	s.sql = b
}

func (s *SqlState) GetSql() []byte {
	return s.sql
}
func (s *SqlState) GetArgs() []interface{} {
	return s.args
}

type IDriver interface {
	AdaptSql([]byte) []byte
	ConnectString() string
	SetConnectString(string)
	QuoteField(string) string
	LastInsertID(string, string) string
	FlatData(reflect.Type, interface{}) interface{}
	StringSlice([]byte) ([]string, error)
	Int64Slice([]byte) ([]int64, error)
	ParseNumberSlice([]byte, interface{}) error
	ParseStringSlice([]byte, interface{}) error
	HStore([]byte) (map[string]string, error)
	DriverName() string
}

type Connection struct {
	Driver IDriver
	// counter int
	conn            *sql.DB
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func NewConnection(driver IDriver) *Connection {
	return &Connection{Driver: driver}
}

func (c *Connection) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Connection) Conn() *sql.DB {
	if c.conn == nil {
		c.conn = c.Connect()
	}

	// if err := c.conn.Ping(); err != nil {
	// 	c.conn.Close()
	// 	c.conn = c.Connect()
	// }
	return c.conn
}

func (c *Connection) Connect() *sql.DB {
	conn, err := sql.Open(c.Driver.DriverName(), c.Driver.ConnectString())
	if err != nil {
		panic(err)
	}

	conn.SetMaxOpenConns(c.MaxOpenConns)
	conn.SetMaxIdleConns(c.MaxIdleConns)
	conn.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second)

	return conn
}

type Database struct {
	Name       string
	Driver     IDriver
	Log        *log.Logger
	Connection *Connection
	Location   *time.Location
}

func NewDatabase(name string, driver IDriver, lg *log.Logger) *Database {
	conn := NewConnection(driver)

	conn.MaxOpenConns = dbConf.DefaultInt("max_open_conns", 100)
	conn.MaxIdleConns = dbConf.DefaultInt("max_idle_conns", 10)
	conn.ConnMaxLifetime = dbConf.DefaultInt("conn_max_life_time", 30)

	return &Database{Name: name, Driver: driver, Log: lg, Connection: conn}
}

func (this *Database) isError(err error) bool {
	if err != nil {
		this.Log.Error("[sql]", err.Error())
		return true
	}
	return false
}

func (this *Database) QueryPrepare(s *SqlState) (DataSet, error) {
	return QueryPrepare(this.Connection.Conn(), s.GetSql(), s.GetArgs()...)
}

// func (this *Database) FindPrepare(cls IVO, s *SqlState) (VODataSet, error) {
// 	return FindPrepare(this.Connection.Conn(), cls, s.GetSql(), s.GetArgs()...)
// }

func (this *Database) Query(s *SqlState) (DataSet, error) {
	return Query(this.Connection.Conn(), s.GetSql(), s.GetArgs()...)
}

// func (this *Database) Find(cls IVO, s *SqlState) (VODataSet, error) {
// 	return Find(this.Connection.Conn(), cls, s.GetSql(), s.GetArgs()...)
// }

func (this *Database) ExecPrepare(s *SqlState) (sql.Result, error) {
	return ExecPrepare(this.Connection.Conn(), s.GetSql(), s.GetArgs()...)
}

func (this *Database) Exec(s *SqlState) (sql.Result, error) {
	return Exec(this.Connection.Conn(), s.GetSql(), s.GetArgs()...)
}

// ExecStr
// Exec by string
func (this *Database) ExecStr(sqlstr string, args ...interface{}) (sql.Result, error) {
	return this.Exec(NewSqlState([]byte(sqlstr), args...))
}
