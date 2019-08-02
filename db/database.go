package db

import (
	"database/sql"
	"time"

	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
)

// IDriver interface
type IDriver interface {
	Adapt(string, int) []byte
	ConnectString() string

	QuoteIdentifier(string) string // indentifier
	QuoteIdentifierB(string) []byte
	QuoteLiteral(string) string // store value

	LastInsertID(string, string) string
	StoreData(key string, val interface{}) interface{}
	Strings([]byte) ([]string, error)
	Int64s([]byte) ([]int64, error)
	Floats([]byte) ([]float64, error)
	Ints([]byte) ([]int, error)
	ParseNumberSlice([]byte, interface{}) error
	ParseStringSlice([]byte, interface{}) error
	// HStore([]byte) (map[string]string, error)
	Name() string
}

//Database class
type Database struct {
	Name   string
	Driver IDriver
	log    *log.Logger

	db              *sql.DB
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// NewDatabase new
func NewDatabase(name string, driver IDriver, dbConf conf.Conf, lg *log.Logger) *Database {
	d := &Database{Name: name, Driver: driver, log: lg}
	d.MaxOpenConns = dbConf.DefaultInt("max_open_conns", 1000)
	d.MaxIdleConns = dbConf.DefaultInt("max_idle_conns", 10)
	d.ConnMaxLifetime = dbConf.DefaultInt("conn_max_life_time", 60)

	return d
}

// Conn DB
func (d *Database) Conn() *sql.DB {
	var err error
	if d.db == nil {
		d.db, err = d.Connect()
		if err != nil {
			return d.db
		}
		return d.db
	}

	err = d.db.Ping()
	if err != nil {
		d.db, err = d.Connect()
		if err != nil {
			return d.db
		}
		return d.db
	}

	return d.db
}

// Connect db
func (d *Database) Connect() (*sql.DB, error) {
	db, err := sql.Open(d.Driver.Name(), d.Driver.ConnectString())
	if err != nil {
		d.log.Crit(err)
		return db, err
	}

	db.SetMaxOpenConns(d.MaxOpenConns)
	db.SetMaxIdleConns(d.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(d.ConnMaxLifetime) * time.Second)

	return db, nil
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
	result, _, err := cQuery(0, 1, d.Conn(), sqlstr, args...)
	return result, err
}

//Query db
func (d *Database) Query(sqlstr string, args ...interface{}) (DataSet, error) {
	d.Log(sqlstr, args)
	dataset, _, err := cQuery(0, 0, d.Conn(), sqlstr, args...)
	return dataset, err
}

// QueryRowsPrepare db
func (d *Database) QueryRowsPrepare(sqlstr string, args ...interface{}) (MapRows, error) {
	d.Log(sqlstr, args)
	_, rows, err := cQuery(1, 1, d.Conn(), sqlstr, args...)
	return rows, err
}

//QueryRows db
func (d *Database) QueryRows(sqlstr string, args ...interface{}) (MapRows, error) {
	d.Log(sqlstr, args)
	_, rows, err := cQuery(1, 0, d.Conn(), sqlstr, args...)
	return rows, err
}

// ExecPrepare db
func (d *Database) ExecPrepare(sqlstr string, args ...interface{}) (sql.Result, error) {
	d.Log(sqlstr, args)
	return ExecPrepare(d.Conn(), sqlstr, args...)
}

//Exec db
func (d *Database) Exec(sqlstr string, args ...interface{}) (sql.Result, error) {
	d.Log(sqlstr, args)
	return Exec(d.Conn(), sqlstr, args...)
}
