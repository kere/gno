package dba

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
)

var (
	dbpool *databasePool
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

func confGet(config map[string]string, key, defaultValue string) string {
	if v, ok := config[key]; ok {
		return v
	}
	return defaultValue
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

	driverName := confGet(c, "driver", "postgres")
	logger := NewLogger(c)

	var driver IDriver
	switch driverName {
	case "postgres", "psql":
		driver = &Postgres{DBName: confGet(c, "dbname", "app"),
			User:     confGet(c, "user", "postgres"),
			Password: confGet(c, "password", "123"),
			Host:     confGet(c, "host", "127.0.0.1"),
			HostAddr: confGet(c, "hostaddr", ""),
			Port:     confGet(c, "port", "5432"),
		}

	// case "mysql":
	// 	driver = &drivers.Mysql{DBName: confGet(c, "dbname"),
	// 		User:       confGet(c, "user"),
	// 		Password:   confGet(c, "password"),
	// 		Protocol:   confGet(c, "protocol"),
	// 		Parameters: confGet(c, "parameters"),
	// 		Addr:       confGet(c, "addr")}
	//
	// case "sqlite3":
	// 	driver = &drivers.Sqlite3{File: confGet(c, "file")}

	default:
		// driver = &drivers.Common{}
		panic("you may need regist a custom driver: db.RegistDriver(Mysql{})")
	}

	// driver.SetConnectString(confGet(c, "connect"))
	d := NewDatabase(name, driver, conf.Conf(c), logger)

	dbpool.SetDatabase(name, d)
	dbpool.SetCurrent(d)
	return d
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

// NewLogger func
func NewLogger(dbConf map[string]string) *log.Logger {
	conf := conf.Conf(dbConf)
	levelStr := conf.Get("level")
	level := log.IntLevel(levelStr)

	name := conf.Get("logname")
	if name == "" {
		name = "db"
	}

	if level < 0 {
		return log.NewEmpty()
	}

	folder := filepath.Join(filepath.Dir(os.Args[0]), "/var/log/")
	return log.New(folder, name, conf.Get("logstore"), levelStr)
}
