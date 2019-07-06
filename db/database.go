package db

import (
	"database/sql"
	"reflect"
	"sync"
	"time"

	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
)

// IDriver interface
type IDriver interface {
	Adapt(string, int) string
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
	Name() string
}

// Connection class
type Connection struct {
	Driver IDriver
	// counter int
	conn            *sql.DB
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	Mutex           sync.Mutex
}

// NewConnection new
func NewConnection(driver IDriver) *Connection {
	return &Connection{Driver: driver}
}

// Close conn
func (c *Connection) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Conn sql.DB
func (c *Connection) Conn() *sql.DB {
	if c.conn == nil {
		c.conn = c.Connect()
	}

	if err := c.conn.Ping(); err != nil {
		c.conn.Close()
		c.conn = c.Connect()
	}
	return c.conn
}

// Connect db
func (c *Connection) Connect() *sql.DB {
	conn, err := sql.Open(c.Driver.Name(), c.Driver.ConnectString())
	if err != nil {
		panic(err)
	}

	conn.SetMaxOpenConns(c.MaxOpenConns)
	conn.SetMaxIdleConns(c.MaxIdleConns)
	conn.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second)

	return conn
}

//Database class
type Database struct {
	Name       string
	Driver     IDriver
	log        *log.Logger
	Connection *Connection
	// Location   *time.Location
}

// NewDatabase new
func NewDatabase(name string, driver IDriver, dbConf conf.Conf, lg *log.Logger) *Database {
	conn := NewConnection(driver)

	conn.MaxOpenConns = dbConf.DefaultInt("max_open_conns", 100)
	conn.MaxIdleConns = dbConf.DefaultInt("max_idle_conns", 10)
	conn.ConnMaxLifetime = dbConf.DefaultInt("conn_max_life_time", 30)

	return &Database{Name: name, Driver: driver, log: lg, Connection: conn}
}

// func (d *Database) isError(err error) bool {
// 	if err != nil {
// 		d.log.Error(d.Name, err.Error())
// 		return true
// 	}
// 	return false
// }

// Log db
func (d *Database) Log(sql string, args []interface{}) {
	d.log.Sql(d.Name, sql, args)
}

// SetLogLevel db
func (d *Database) SetLogLevel(level string) {
	d.log.SetLevel(level)
}

// SetLog db
func (d *Database) SetLog(l *log.Logger) {
	d.log = l
}

// QueryPrepare db
func (d *Database) QueryPrepare(sql string, args ...interface{}) (DataSet, error) {
	d.Log(sql, args)
	return CQueryPrepare(d.Connection.Conn(), sql, args...)
}

//Query db
func (d *Database) Query(sql string, args ...interface{}) (DataSet, error) {
	d.Log(sql, args)
	return CQuery(d.Connection.Conn(), sql, args...)
}

// ExecPrepare db
func (d *Database) ExecPrepare(sql string, args ...interface{}) (sql.Result, error) {
	d.Log(sql, args)
	return ExecPrepare(d.Connection.Conn(), sql, args...)
}

//Exec db
func (d *Database) Exec(sql string, args ...interface{}) (sql.Result, error) {
	d.Log(sql, args)
	return Exec(d.Connection.Conn(), sql, args...)
}
