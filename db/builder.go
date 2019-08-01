package db

import (
	"database/sql"
	"reflect"
)

var ivotype = reflect.TypeOf((*IVO)(nil)).Elem()

type builder struct {
	conn     *sql.DB
	database *Database
}

func (b *builder) GetDatabase() *Database {
	if b.database != nil {
		return b.database
	}
	b.database = Current()
	return b.database
}

func (b *builder) SetDatabase(d *Database) {
	b.database = d
}
