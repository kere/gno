package dba

import (
	"database/sql"
	"io"
	"time"

	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
)

// IDriver interface
type IDriver interface {
	Name() string
	ConnectString() string

	WriteQuoteIdentifier(io.Writer, string)

	LastInsertID(string, string) string
	StoreData(key string, val interface{}) interface{}

	Strings([]byte) ([]string, error)
	Int64s([]byte) ([]int64, error)
	Floats([]byte) ([]float64, error)
	Ints([]byte) ([]int, error)
	ParseNumberSlice([]byte, interface{}) error
	ParseStringSlice([]byte, interface{}) error
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
	d.ConnMaxLifetime = dbConf.DefaultInt("conn_max_life_time", 30)

	return d
}

// Conn DB
func (d *Database) DB() *sql.DB {
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
