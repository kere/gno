package db

// to store value

import (
	"bytes"
	"fmt"

	"github.com/valyala/bytebufferpool"
)

func sqlUpdateParamsByMapRow(row MapRow) ([]byte, []interface{}) {
	l := len(row)
	keys := make([][]byte, l)
	seq := 1
	i := 0

	database := Current()
	values := make([]interface{}, l)

	for k := range row {
		if row[k] == nil {
			// value=NULL
			tmp := append(database.Driver.QuoteIdentifierB(k), '=')
			keys[i] = append(tmp, BNull...)
		} else {
			// value != nil
			tmp := append(database.Driver.QuoteIdentifierB(k), '=', '$')
			keys[i] = append(tmp, []byte(fmt.Sprint(seq))...)
		}

		values[i] = database.Driver.StoreData(k, row[k])

		i++
		seq++
	}
	return bytes.Join(keys, BComma), values
}

func sqlInsertParamsByMapRow(row MapRow) ([]byte, []interface{}) {
	l := len(row)
	keys := make([][]byte, l)
	i := 0

	database := Current()
	values := make([]interface{}, l)

	for k := range row {
		keys[i] = database.Driver.QuoteIdentifierB(k)
		values[i] = database.Driver.StoreData(k, row[k])
		i++
	}
	return bytes.Join(keys, BComma), values
}

func sqlInsertParams(fields []string, row Row) ([]byte, []interface{}) {
	n := len(fields)
	keys := make([][]byte, n)
	database := Current()
	values := make([]interface{}, n)

	for i := 0; i < n; i++ {
		keys[i] = database.Driver.QuoteIdentifierB(fields[i])
		values[i] = database.Driver.StoreData(fields[i], row[i])
		i++
	}
	return bytes.Join(keys, BComma), values
}

func sqlInsertKeysByMapRow(row MapRow) ([]byte, []string) {
	bkeys := make([][]byte, len(row))
	keys := make([]string, len(row))
	i := 0

	database := Current()

	for k := range row {
		keys[i] = k
		bkeys[i] = database.Driver.QuoteIdentifierB(k)
		i++
	}
	return bytes.Join(bkeys, BComma), keys
}

func writeInsertMByMapRow(buf *bytebufferpool.ByteBuffer, keys []string, rows MapRows) []interface{} {
	l := len(rows)
	n := len(keys)
	seq := 1

	database := Current()
	values := make([]interface{}, 0, l*n+10)

	for i := 0; i < l; i++ {
		row := rows[i]
		buf.WriteByte('(')

		for k := 0; k < n; k++ {
			key := keys[k]
			val := database.Driver.StoreData(key, row[key])

			values = append(values, val)

			if seq == 1 {
				buf.WriteString("$1")
			} else {
				buf.WriteString(fmt.Sprint(SDoller, seq))
			}
			if k < n-1 {
				buf.WriteByte(',')
			}
			seq++
		}

		buf.WriteByte(')')
		if i < l-1 {
			buf.WriteByte(',')
		}

	}
	buf.WriteByte(';')

	return values
}

func writeInsertMByDataSet(buf *bytebufferpool.ByteBuffer, dataset *DataSet) []interface{} {
	l := dataset.Len()
	n := len(dataset.Fields)
	seq := 1

	database := Current()
	values := make([]interface{}, 0, l*n+10)
	cols := dataset.Columns

	for i := 0; i < l; i++ {
		buf.WriteByte('(')

		for k := 0; k < n; k++ {
			key := dataset.Fields[k]
			v := cols[k][i]
			val := database.Driver.StoreData(key, v)

			values = append(values, val)

			if seq == 1 {
				buf.WriteString("$1")
			} else {
				buf.WriteString(fmt.Sprint(SDoller, seq))
			}
			if k < n-1 {
				buf.WriteByte(',')
			}
			seq++
		}

		buf.WriteByte(')')
		if i < l-1 {
			buf.WriteByte(',')
		}

	}
	buf.WriteByte(';')

	return values
}

// func sqlInsertMByMapRow(keys []string, row MapRow, values []interface{}) []byte {
// 	l := len(keys)
// 	database := Current()
// 	sbu := strings.Builder{}
//
// 	for i := 0; i < l; i++ {
// 		k := keys[i]
//
// 		var val interface{}
// 		if len(k) > 5 && k[len(k)-5:] == subfixJSON {
// 			b, _ := json.Marshal(row[k])
// 			val = string(b)
// 		} else {
// 			val = database.Driver.StoreData(reflect.TypeOf(row[k]), row[k])
// 		}
//
// 		switch val.(type) {
// 		case time.Time:
// 			values[i] = SQuot + (row[k].(time.Time)).Format(time.RFC1123) + SQuot
// 		case string:
// 			values[i] = SQuot + row[k].(string) + SQuot
// 		case []byte:
// 			values[i] = SQuot + string(row[k].([]byte)) + SQuot
// 		default:
// 			values[i] = fmt.Sprint(row[k])
// 		}
// 		i++
// 	}
//
// 	return strings.Join(values, SCommaSplit)
// }

// // keyValueList
// // stype:insert,update
// func keyValueList(actionType string, data interface{}) (keys [][]byte, values []interface{}, stmts [][]byte) {
// 	var d map[string]interface{}
// 	switch data.(type) {
// 	case MapRow:
// 		d = map[string]interface{}(data.(MapRow))
//
// 	case map[string]interface{}:
// 		d = data.(map[string]interface{})
//
// 	default:
// 		sm := NewStructConvert(data)
// 		d = map[string]interface{}(sm.Struct2DataRow(actionType))
//
// 	}
//
// 	l := len(d)
// 	isUpdate := actionType == ActionUpdate
// 	keys = make([][]byte, l)
// 	values = make([]interface{}, 0)
// 	stmts = make([][]byte, l)
// 	database := Current()
// 	i, ii := 0, 1
// 	for k, v := range d {
// 		typ := reflect.TypeOf(v)
//
// 		if isUpdate {
// 			if v == nil {
// 				// version=version+1
// 				if strings.IndexByte(k, BEqual[0]) > 0 {
// 					keys[i] = []byte(k)
// 				} else {
// 					// field=NULL
// 					arr := append([]byte(database.Driver.QuoteField(k)), '=')
// 					keys[i] = append(arr, BNull...)
// 				}
// 				i++
// 				continue
// 			} else {
// 				arr := append([]byte(database.Driver.QuoteField(k)), '=', '$')
// 				keys[i] = append(arr, []byte(fmt.Sprint(ii))...)
// 			}
// 		} else {
// 			keys[i] = []byte(database.Driver.QuoteField(k))
// 		}
//
// 		if !isUpdate && v == nil {
// 			stmts[i] = BNull //insert 时为null
// 		} else {
// 			stmts[i] = append([]byte("$"), []byte(fmt.Sprint(ii))...)
// 			ii++
// 		}
//
// 		if v != nil && typ.Kind() == reflect.Ptr {
// 			typ = typ.Elem()
// 		}
//
// 		if v != nil && typ.Implements(ivotype) {
// 			values = append(values, NewStructConvert(v).Struct2DataRow(actionType))
// 		} else if len(k) > 5 && k[len(k)-5:] == subfixJSON {
// 			b, _ := json.Marshal(v)
// 			values = append(values, b)
// 		} else {
// 			values = append(values, database.Driver.StoreData(typ, v))
// 		}
//
// 		i++
// 	}
//
// 	return
// }
