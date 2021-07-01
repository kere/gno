package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
)

var (
	bytePool bytebufferpool.Pool
)

//Postgres class
type Postgres struct {
	DBName   string
	User     string
	Password string
	Host     string
	HostAddr string
	Port     string
}

// Name f
func (p *Postgres) Name() string {
	return DriverPSQL
}

// // Adapt f
// func (p *Postgres) Adapt(sqlstr string, n int) string {
// 	src := util.Str2Bytes(sqlstr)
// 	arr := bytes.Split(src, BQuestionMark)
// 	// arr := strings.Split(sqlstr, sQuestionMark)
// 	l := len(arr)
// 	if l == 0 {
// 		return ""
// 	}
// 	if l == 1 {
// 		return sqlstr
// 	}
//
// 	buf := bytePool.Get()
// 	for i := 0; i < l-1; i++ {
// 		buf.Write(arr[i])
// 		buf.WriteByte('$')
//
// 		buf.WriteString(fmt.Sprint(i + 1 + n))
// 	}
//
// 	b := buf.String()
// 	bytePool.Put(buf)
//
// 	return b
// }

// ConnectString f
func (p *Postgres) ConnectString() string {
	if p.Host == "" {
		p.Host = "127.0.0.1"
	}
	if p.Port == "" {
		p.Port = "5432"
	}

	if p.HostAddr != "" {
		return fmt.Sprintf("dbname=%s user=%s password=%s hostaddr=%s sslmode=disable",
			p.DBName,
			p.User,
			p.Password,
			p.HostAddr)

	}
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		p.DBName,
		p.User,
		p.Password,
		p.Host,
		p.Port)
}

// LastInsertID f
func (p *Postgres) LastInsertID(table, pkey string) string {
	// return "select currval(pg_get_serial_sequence('" + table + "','" + pkey + "'))"
	return fmt.Sprint("SELECT currval(pg_get_serial_sequence('", table, "','", pkey, "')) as count")
}

func sliceToStore(v interface{}) string {
	val := reflect.ValueOf(v)
	l := val.Len()
	if l == 0 {
		return "{}"
	}
	arr := make([]string, 0, l)

	for i := 0; i < l; i++ {
		v := fmt.Sprint(val.Index(i).Interface())
		if v == "" {
			continue
		}
		arr = append(arr, v)
	}
	return fmt.Sprint("{", strings.Join(arr, ","), "}")
}

// StoreData for value
func (p *Postgres) StoreData(key string, v interface{}) interface{} {
	if v == nil {
		return nil
	}
	if strings.HasSuffix(key, "_json") {
		src, err := json.Marshal(v)
		if err != nil {
			return util.BNull
		}
		return src
	}

	typ := reflect.TypeOf(v)
	switch typ.Kind() {
	default:
		return v

	case reflect.Array:
		return sliceToStore(v)

	case reflect.Slice:
		switch v.(type) {
		case []byte:
			return v
		default:
			return sliceToStore(v)
		}

	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			return v

		default:
			b, _ := json.Marshal(v)
			return b
		}
	}
}

// Strings []string
func (p *Postgres) Strings(src []byte) ([]string, error) {
	return doStrings(src, true)
}

// StringsNotSafe []string
func (p *Postgres) StringsNotSafe(src []byte) ([]string, error) {
	return doStrings(src, false)
}

// BytesArr [][]byte
func (p *Postgres) BytesArr(src []byte) ([][]byte, error) {
	return doByteArr(src, true)
}

// BytesArrNotSafe [][]byte
func (p *Postgres) BytesArrNotSafe(src []byte) ([][]byte, error) {
	return doByteArr(src, false)
}

// Int64s arr
func (p *Postgres) Int64s(src []byte) ([]int64, error) {
	return doInt64s(src, false)
}

// Int64sP arr
func (p *Postgres) Int64sP(src []byte) ([]int64, error) {
	return doInt64s(src, true)
}

// Floats arr
func (p *Postgres) Floats(src []byte) ([]float64, error) {
	return doFloats(src, false)
}

// FloatsP arr
func (p *Postgres) FloatsP(src []byte) ([]float64, error) {
	return doFloats(src, true)
}

// Ints arr
func (p *Postgres) Ints(src []byte) ([]int, error) {
	return doInts(src, false)
}

// Ints arr
func (p *Postgres) IntsP(src []byte) ([]int, error) {
	return doInts(src, true)
}

func doInt64s(src []byte, isPool bool) ([]int64, error) {
	if len(src) < 2 {
		return nil, nil
	}
	b := src

	if bytes.HasPrefix(src, util.BBraceLeft) {
		b = src[1:]
	}
	if bytes.HasSuffix(b, util.BBraceRight) {
		b = b[:len(b)-1]
	}

	if isPool {
		return util.SplitBytes2Int64P(b, util.BComma)
	}
	return util.SplitBytes2Int64(b, util.BComma)
	// var result []int64
	// if isPool {
	// 	result = util.GetInt64s()
	// } else {
	// 	result = make([]int64, 0, 20)
	// }
	//
	// count := bytes.Count(b, util.BComma)
	// var v int64
	// var err error
	// var index int
	// for i := 0; i < count; i++ {
	// 	index = bytes.Index(b, util.BComma)
	// 	if index == -1 {
	// 		break
	// 	}
	//
	// 	switch util.BytesNumType(b[:index]) {
	// 	case 'f':
	// 		var val float64
	// 		val, err = strconv.ParseFloat(util.Bytes2Str(b[:index]), 64)
	// 		v = int64(val)
	// 	case 'i':
	// 		v, err = strconv.ParseInt(util.Bytes2Str(b[:index]), 10, 64)
	// 	default:
	// 		return result, errors.New("do Ints:can not to parse str to num")
	// 	}
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result = append(result, v)
	// 	if index == len(b)-1 {
	// 		b = nil
	// 	} else {
	// 		b = b[index+1:]
	// 	}
	// }
	//
	// if len(b) > 0 {
	// 	switch util.BytesNumType(b) {
	// 	case 'f':
	// 		var val float64
	// 		val, err = strconv.ParseFloat(util.Bytes2Str(b), 64)
	// 		v = int64(val)
	// 	case 'i':
	// 		v, err = strconv.ParseInt(util.Bytes2Str(b), 10, 64)
	// 	default:
	// 		return result, errors.New("do Ints:can not to parse str to num")
	// 	}
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result = append(result, v)
	// }
	// return result, nil
}

func doInts(src []byte, isPool bool) ([]int, error) {
	if len(src) < 2 {
		return nil, nil
	}
	b := src

	if bytes.HasPrefix(src, util.BBraceLeft) {
		b = src[1:]
	}
	if bytes.HasSuffix(b, util.BBraceRight) {
		b = b[:len(b)-1]
	}

	if isPool {
		return util.SplitBytes2IntP(b, util.BComma)
	}
	return util.SplitBytes2Int(b, util.BComma)

	// var result []int
	// if isPool {
	// 	result = util.GetInts()
	// } else {
	// 	result = make([]int, 0, 20)
	// }
	//
	// count := bytes.Count(b, util.BComma)
	// var index int
	// var v int64
	// var err error
	// for i := 0; i < count; i++ {
	// 	index = bytes.Index(b, util.BComma)
	// 	if index == -1 {
	// 		break
	// 	}
	//
	// 	switch util.BytesNumType(b[:index]) {
	// 	case 'f':
	// 		var val float64
	// 		val, err = strconv.ParseFloat(util.Bytes2Str(b[:index]), 64)
	// 		v = int64(val)
	// 	case 'i':
	// 		v, err = strconv.ParseInt(util.Bytes2Str(b[:index]), 10, 64)
	// 	default:
	// 		return result, errors.New("do Ints:can not to parse str to num")
	// 	}
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result = append(result, int(v))
	//
	// 	if index == len(b)-1 {
	// 		b = nil
	// 	} else {
	// 		b = b[index+1:]
	// 	}
	// }
	// if len(b) > 0 {
	// 	switch util.BytesNumType(b) {
	// 	case 'f':
	// 		var val float64
	// 		val, err = strconv.ParseFloat(util.Bytes2Str(b), 64)
	// 		v = int64(val)
	// 	case 'i':
	// 		v, err = strconv.ParseInt(util.Bytes2Str(b), 10, 64)
	// 	default:
	// 		return result, errors.New("do Ints:can not to parse str to num")
	// 	}
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result = append(result, int(v))
	// }
	// return result, nil
}

func doFloats(src []byte, isPool bool) ([]float64, error) {
	if len(src) < 2 {
		return nil, nil
	}
	b := src

	if bytes.HasPrefix(src, util.BBraceLeft) {
		b = src[1:]
	}
	if bytes.HasSuffix(b, util.BBraceRight) {
		b = b[:len(b)-1]
	}

	if isPool {
		return util.SplitBytes2FloatsP(b, util.BComma)
	}
	return util.SplitBytes2Floats(b, util.BComma)
	// var result []float64
	// if isPool {
	// 	result = util.GetFloats()
	// } else {
	// 	result = make([]float64, 0, 20)
	// }
	//
	// count := bytes.Count(b, util.BComma)
	// var index int
	// for i := 0; i < count; i++ {
	// 	index = bytes.Index(b, util.BComma)
	// 	if index == -1 {
	// 		break
	// 	}
	// 	v, err := strconv.ParseFloat(util.Bytes2Str(b[:index]), 64)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result = append(result, v)
	//
	// 	if index == len(b)-1 {
	// 		b = nil
	// 	} else {
	// 		b = b[index+1:]
	// 	}
	// }
	// if len(b) > 0 {
	// 	v, err := strconv.ParseFloat(util.Bytes2Str(b), 64)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result = append(result, v)
	// }
	// return result, nil
}
func doStrings(src []byte, isSafe bool) ([]string, error) {
	if len(src) < 2 {
		return nil, nil
	}
	b := src

	if bytes.HasPrefix(src, util.BBraceLeft) {
		b = src[1:]
	}
	if bytes.HasSuffix(b, util.BBraceRight) {
		b = b[:len(b)-1]
	}

	if isSafe {
		arr := util.SplitBytesNotSafe(b, util.BComma)
		count := len(arr)
		result := make([]string, count)
		for i := 0; i < count; i++ {
			result[i] = string(arr[i])
		}
		return result, nil
	}
	return util.SplitStrNotSafe(util.Bytes2Str(b), util.SComma), nil
}
func doByteArr(src []byte, isSafe bool) ([][]byte, error) {
	if len(src) < 2 {
		return nil, nil
	}
	b := src

	if bytes.HasPrefix(src, util.BBraceLeft) {
		b = src[1:]
	}
	if bytes.HasSuffix(b, util.BBraceRight) {
		b = b[:len(b)-1]
	}

	if isSafe {
		arr := util.SplitBytesNotSafe(b, util.BComma)
		count := len(arr)
		for i := 0; i < count; i++ {
			row := make([]byte, len(arr[i]))
			copy(row, arr[i])
			arr[i] = row
		}
		return arr, nil
	}
	return util.SplitBytesNotSafe(b, util.BComma), nil
}

// WriteQuoteIdentifier f
func (p *Postgres) WriteQuoteIdentifier(w io.Writer, s string) {
	// return pq.QuoteLiteral(literal)
	// `"` + strings.Replace(name, `"`, `""`, -1) + `"`
	str := strings.Replace(s, `"`, `""`, -1)
	w.Write(util.BDoubleQuote)
	w.Write(util.Str2Bytes(str))
	w.Write(util.BDoubleQuote)
}

// // QuoteLiteral quotes a 'literal' (e.g. a parameter, often used to pass literal
// // to DDL and other statements that do not accept parameters) to be used as part
// // of an SQL statement.  For example:
// //
// //    exp_date := pq.QuoteLiteral("2023-01-05 15:00:00Z")
// //    err := db.Exec(fmt.Sprintf("CREATE ROLE my_user VALID UNTIL %s", exp_date))
// //
// // Any single quotes in name will be escaped. Any backslashes (i.e. "\") will be
// // replaced by two backslashes (i.e. "\\") and the C-style escape identifier
// // that PostgreSQL provides ('E') will be prepended to the string.
// func (p *Postgres) QuoteLiteral(literal string) string {
// 	return pq.QuoteLiteral(literal)
// }
