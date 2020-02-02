package util

import "sync"

var (
	rowPool    sync.Pool
	colPool    sync.Pool
	intsPool   sync.Pool
	strsPool   sync.Pool
	int64sPool sync.Pool
)

func parseArgs(args []int) (int, int) {
	l := len(args)
	switch {
	case l == 0:
		return 0, 0
	case l == 1:
		return args[0], args[0]
	default:
		return args[0], args[1]
	}
}

// GetRow from pool
func GetRow(args ...int) []float64 {
	l, capN := parseArgs(args)
	v := rowPool.Get()
	if capN < 20 {
		capN = 20
	}

	if v == nil {
		return make([]float64, l, capN)
	}
	row := v.([]float64)

	for i := 0; i < l; i++ {
		row = append(row, 0)
	}
	return row
}

// PutRow to pool
func PutRow(r []float64) {
	if cap(r) == 0 {
		return
	}
	rowPool.Put(r[:0])
}

// ---------- float64 ----------

// GetColumn from pool
func GetColumn(args ...int) []float64 {
	l, capN := parseArgs(args)
	v := colPool.Get()
	if v == nil {
		if capN < 100 {
			capN = 100
		}
		return make([]float64, l, capN)
	}
	col := v.([]float64)

	for i := 0; i < l; i++ {
		col = append(col, 0)
	}
	return col
}

// PutColumn to pool
func PutColumn(r []float64) {
	if cap(r) == 0 {
		return
	}
	colPool.Put(r[:0])
}

// ---------- int64 ----------

// GetInt64 from pool
func GetInt64(args ...int) []int64 {
	l, capN := parseArgs(args)
	v := int64sPool.Get()
	if v == nil {
		if capN < 50 {
			capN = 50
		}
		return make([]int64, l, capN)
	}

	arr := (v.([]int64))
	for i := 0; i < l; i++ {
		arr = append(arr, 0)
	}
	return arr
}

// PutInt64 to pool
func PutInt64(r []int64) {
	if cap(r) == 0 {
		return
	}
	int64sPool.Put(r[:0])
}

// ---------- int ----------

// GetInt from pool, with 0 value
func GetInt(args ...int) []int {
	l, capN := parseArgs(args)
	v := intsPool.Get()
	if v == nil {
		if capN < 50 {
			capN = 50
		}
		return make([]int, l, capN)
	}
	arr := v.([]int)

	for i := 0; i < l; i++ {
		arr = append(arr, 0)
	}
	return arr
}

// PutInt to pool
func PutInt(r []int) {
	if cap(r) == 0 {
		return
	}
	intsPool.Put(r[:0])
}

// ---------- strings ----------

// GetStrings from pool
func GetStrings(args ...int) []string {
	l, capN := parseArgs(args)
	v := strsPool.Get()
	if v == nil {
		if capN < 50 {
			capN = 50
		}
		return make([]string, l, capN)
	}
	arr := v.([]string)

	for i := 0; i < l; i++ {
		arr = append(arr, "")
	}
	return arr
}

// PutStrings to pool
func PutStrings(r []string) {
	if cap(r) == 0 {
		return
	}
	strsPool.Put(r[:0])
}
