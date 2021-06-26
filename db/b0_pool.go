package db

import (
	"sync"

	"github.com/kere/gno/libs/util"
)

var (
	rowPool sync.Pool
	colPool sync.Pool
	lockA   sync.Mutex
	lockCol sync.Mutex
)

// GetDataSet get from pool
// l 初始化的长度
func GetDataSet(fields []string, args ...int) DataSet {
	lockA.Lock()
	n := len(fields)
	l, capN := util.ParsePoolArgs(args, 100)

	cols := make([][]interface{}, n)
	for i := 0; i < n; i++ {
		cols[i] = GetColumn(l, capN)
	}

	lockA.Unlock()
	return DataSet{Fields: fields, Columns: cols}
}

// PutDataSet series into pool
func PutDataSet(dat *DataSet) {
	if dat == nil || cap(dat.Columns) == 0 {
		dat.Release()
		return
	}
	lockA.Lock()

	n := len(dat.Columns)
	for i := 0; i < n; i++ {
		PutColumn(dat.Columns[i])
	}
	dat.Release()
	lockA.Unlock()
}

// GetRow from pool
func GetRow(args ...int) []interface{} {
	l, capN := util.ParsePoolArgs(args, 20)
	v := rowPool.Get()
	if v == nil {
		return make([]interface{}, l, capN)
	}
	row := v.([]interface{})

	for i := 0; i < l; i++ {
		row = append(row, 0)
	}
	return row
}

// PutRow to pool
func PutRow(row []interface{}) {
	if cap(row) == 0 {
		return
	}
	rowPool.Put(row[:0])
}

// GetColumn from pool
func GetColumn(args ...int) []interface{} {
	lockCol.Lock()
	defer lockCol.Unlock()
	l, capN := util.ParsePoolArgs(args, 100)
	v := colPool.Get()
	if v == nil {
		return make([]interface{}, l, capN)
	}
	col := v.([]interface{})

	for i := 0; i < l; i++ {
		col = append(col, 0)
	}
	return col
}

// PutColumn to pool
func PutColumn(r []interface{}) {
	lockCol.Lock()
	defer lockCol.Unlock()
	if cap(r) == 0 {
		return
	}
	colPool.Put(r[:0])
}
