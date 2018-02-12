package db

import (
	"database/sql"
	"reflect"
)

var ivotype reflect.Type = reflect.TypeOf((*IVO)(nil)).Elem()

type builder struct {
	conn *sql.DB
}

func (b *builder) getDatabase() *Database {
	return Current()
}

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
	values = make([]interface{}, l)
	stmts = make([][]byte, l)
	database := Current()
	i := 0
	for k, v := range d {
		keys[i] = []byte(database.Driver.QuoteField(k))

		if isUpdate {
			keys[i] = append(keys[i], B_Equal[0], B_QuestionMark[0])
		}
		stmts[i] = B_QuestionMark

		typ := reflect.TypeOf(v)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		if typ.Implements(ivotype) {
			sm := NewStructConvert(v)
			values[i] = sm.Struct2DataRow(actionType)
		} else {
			values[i] = database.Driver.FlatData(typ, v)
		}

		i++
	}

	return
}
