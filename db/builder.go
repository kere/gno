package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const (
	subfixJSON = "_json"
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

// keyValueList
// stype:insert,update
func keyValueList(actionType string, data interface{}) (keys [][]byte, values []interface{}, stmts [][]byte) {
	var d map[string]interface{}
	switch data.(type) {
	case DataRow:
		d = map[string]interface{}(data.(DataRow))

	case map[string]interface{}:
		d = data.(map[string]interface{})

	default:
		sm := NewStructConvert(data)
		d = map[string]interface{}(sm.Struct2DataRow(actionType))

	}

	l := len(d)
	isUpdate := actionType == ActionUpdate
	keys = make([][]byte, l)
	values = make([]interface{}, 0)
	stmts = make([][]byte, l)
	database := Current()
	i, ii := 0, 1
	for k, v := range d {
		typ := reflect.TypeOf(v)

		if isUpdate {
			if v == nil {
				// version=version+1
				if strings.IndexByte(k, BEqual[0]) > 0 {
					keys[i] = []byte(k)
				} else {
					// field=NULL
					arr := append([]byte(database.Driver.QuoteField(k)), '=')
					keys[i] = append(arr, BNull...)
				}
				i++
				continue
			} else {
				arr := append([]byte(database.Driver.QuoteField(k)), '=', '$')
				keys[i] = append(arr, []byte(fmt.Sprint(ii))...)
			}
		} else {
			keys[i] = []byte(database.Driver.QuoteField(k))
		}

		if !isUpdate && v == nil {
			stmts[i] = BNull //insert 时为null
		} else {
			stmts[i] = append([]byte("$"), []byte(fmt.Sprint(ii))...)
			ii++
		}

		if v != nil && typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		if v != nil && typ.Implements(ivotype) {
			values = append(values, NewStructConvert(v).Struct2DataRow(actionType))
		} else if len(k) > 5 && k[len(k)-5:] == subfixJSON {
			b, _ := json.Marshal(v)
			values = append(values, b)
		} else {
			values = append(values, database.Driver.FlatData(typ, v))
		}

		i++
	}

	return
}
