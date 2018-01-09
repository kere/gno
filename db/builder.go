package db

import "database/sql"

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
func keyValueList(stype string, data interface{}, excludes []string) (keys [][]byte, values []interface{}, stmts [][]byte) {
	var d map[string]interface{}
	switch data.(type) {
	case DataRow:
		d = map[string]interface{}(data.(DataRow))

	case map[string]interface{}:
		d = data.(map[string]interface{})

	default:
		sm := NewStructConvert(data)
		sm.SetExcludes(excludes)
		d = map[string]interface{}(sm.Struct2DataRow(stype))

	}

	l := len(d) - len(excludes)
	isUpdate := stype == "update"
	keys = make([][]byte, l)
	values = make([]interface{}, l)
	stmts = make([][]byte, l)
	database := Current()
	i := 0
	for k, v := range d {
		if InStrings(excludes, k) {
			continue
		}

		keys[i] = []byte(database.Driver.QuoteField(k))

		if isUpdate {
			keys[i] = append(keys[i], B_Equal[0], B_QuestionMark[0])
		}
		stmts[i] = B_QuestionMark

		values[i] = database.Driver.FlatData(v)

		i++
	}

	return
}
