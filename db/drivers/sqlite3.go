package drivers

import "reflect"

type Sqlite3 struct {
	File string
}

func (s *Sqlite3) ConnectString() string {
	return s.File
}

func (s *Sqlite3) QuoteField(str string) string {
	return `"` + str + `"`
}

func (s *Sqlite3) QuoteFieldB(str string) []byte {
	l := len(str)

	// l+8 : "name"=$1234
	arr := make([]byte, l+2, l+8)
	arr[0] = '"'
	for i := 0; i < l; i++ {
		arr[i+1] = str[i]
	}
	arr[l+1] = '"'

	return arr
}

func (s *Sqlite3) LastInsertID(table, id string) string {
	return "SELECT last_insert_rowid()"
}

func (s *Sqlite3) Name() string {
	return DriverSqlite
}

func (s *Sqlite3) Int64Slice(str []byte) ([]int64, error) {
	return nil, nil
}

func (s *Sqlite3) ParseStringSlice(src []byte, ptr interface{}) error {
	return nil
}

func (s *Sqlite3) StoreData(typ reflect.Type, v interface{}) interface{} {
	return nil
}

func (s *Sqlite3) ParseNumberSlice(src []byte, ptr interface{}) error {
	return nil
}

func (s *Sqlite3) StringSlice(src []byte) ([]string, error) {
	return nil, nil
}

func (s *Sqlite3) HStore(src []byte) (map[string]string, error) {
	return nil, nil
}
