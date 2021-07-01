package db

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
	StringsNotSafe([]byte) ([]string, error)
	BytesArr([]byte) ([][]byte, error)
	BytesArrNotSafe([]byte) ([][]byte, error)

	Int64s([]byte) ([]int64, error)
	Int64sP([]byte) ([]int64, error)
	Floats([]byte) ([]float64, error)
	FloatsP([]byte) ([]float64, error)
	Ints([]byte) ([]int, error)
	IntsP([]byte) ([]int, error)
}

// Database class
type Database struct {
	Name   string
	Driver IDriver
	log    *log.Logger

	db *sql.DB

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// NewDatabase new
func NewDatabase(name string, driver IDriver, dbConf conf.Conf, lg *log.Logger) *Database {
	d := &Database{Name: name, Driver: driver, log: lg}
	d.MaxOpenConns = dbConf.DefaultInt("max_open_conns", 300)
	d.MaxIdleConns = dbConf.DefaultInt("max_idle_conns", 50)
	d.ConnMaxLifetime = dbConf.DefaultInt("conn_max_life_time", 30)

	return d
}

// NewBuilder
func (d *Database) NewBuilder(table string) Builder {
	q := Builder{table: table}
	q.database = d
	return q
}

// NewQuery
func (d *Database) NewQuery(table string) QueryBuilder {
	q := NewQuery(table)
	q.database = d
	return q
}

// NewInsert
func (d *Database) NewInsert(table string) InsertBuilder {
	ins := NewInsert(table)
	ins.database = d
	return ins
}

// NewUpdate
func (d *Database) NewUpdate(table string) UpdateBuilder {
	u := NewUpdate(table)
	u.database = d
	return u
}

// NewDelete
func (d *Database) NewDelete(table string) DeleteBuilder {
	del := NewDelete(table)
	del.database = d
	return del
}

// NewExists
func (d *Database) NewExists(table string) ExistsBuilder {
	q := NewExists(table)
	q.database = d
	return q
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
