package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

const (
	subfixJSON = "_json"
)

var ivotype reflect.Type = reflect.TypeOf((*IVO)(nil)).Elem()

type builder struct {
	conn *sql.DB
}

func (b *builder) getDatabase() *Database {
	return Current()
}

// CondParams condition params
type CondParams struct {
	Cond string
	Args []interface{}
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
	isUpdate := actionType == "update"
	keys = make([][]byte, l)
	values = make([]interface{}, 0)
	stmts = make([][]byte, l)
	database := Current()
	i := 0
	for k, v := range d {
		typ := reflect.TypeOf(v)

		keys[i] = []byte(database.Driver.QuoteField(k))

		if isUpdate {
			if v == nil {
				keys[i] = append(keys[i], B_Equal[0])
				keys[i] = append(keys[i], BNull...)
				i++
				continue
			} else {
				keys[i] = append(keys[i], B_Equal[0], B_QuestionMark[0])
			}
		}

		if !isUpdate && v == nil {
			stmts[i] = BNull //insert 时为null
		} else {
			stmts[i] = B_QuestionMark
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
