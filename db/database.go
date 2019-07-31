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
	QuoteFieldB(string) []byte
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
	// if c.conn == nil {
	c.conn = c.Connect()
	// }

	// if err := c.conn.Ping(); err != nil {
	// 	c.conn.Close()
	// 	c.conn = c.Connect()
	// }
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
func (d *Database) QueryPrepare(sqlstr string, args ...interface{}) (DataSet, error) {
	d.Log(sqlstr, args)
	result, _, err := cQuery(0, 1, d.Connection.Conn(), sqlstr, args...)
	return result, err
}

//Query db
func (d *Database) Query(sqlstr string, args ...interface{}) (DataSet, error) {
	d.Log(sqlstr, args)
	dataset, _, err := cQuery(0, 0, d.Connection.Conn(), sqlstr, args...)
	return dataset, err
}

// QueryRowsPrepare db
func (d *Database) QueryRowsPrepare(sqlstr string, args ...interface{}) (MapRows, error) {
	d.Log(sqlstr, args)
	_, rows, err := cQuery(1, 1, d.Connection.Conn(), sqlstr, args...)
	return rows, err
}

//QueryRows db
func (d *Database) QueryRows(sqlstr string, args ...interface{}) (MapRows, error) {
	d.Log(sqlstr, args)
	_, rows, err := cQuery(1, 0, d.Connection.Conn(), sqlstr, args...)
	return rows, err
}

// ExecPrepare db
func (d *Database) ExecPrepare(sqlstr string, args ...interface{}) (sql.Result, error) {
	d.Log(sqlstr, args)
	return ExecPrepare(d.Connection.Conn(), sqlstr, args...)
}

//Exec db
func (d *Database) Exec(sqlstr string, args ...interface{}) (sql.Result, error) {
	d.Log(sqlstr, args)
	return Exec(d.Connection.Conn(), sqlstr, args...)
}
