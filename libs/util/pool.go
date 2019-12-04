package util

import "sync"

var (
	rowPool    sync.Pool
	colPool    sync.Pool
	intsPool   sync.Pool
	strsPool   sync.Pool
	int64sPool sync.Pool
)

// GetRow from pool
func GetRow() []float64 {
	v := rowPool.Get()
	if v == nil {
		return make([]float64, 0, 100)
	}
	return v.([]float64)
}

// GetRowN from pool, with 0 value
func GetRowN(n int) []float64 {
	v := rowPool.Get()
	if v == nil {
		return make([]float64, n)
	}
	row := v.([]float64)

	for i := 0; i < n; i++ {
		row = append(row, 0)
	}
	return row
}

// PutRow to pool
func PutRow(r []float64) {
	rowPool.Put(r[:0])
}

// ---------- float64 ----------

// GetColumn from pool
func GetColumn() []float64 {
	v := colPool.Get()
	if v == nil {
		return make([]float64, 0, 100)
	}
	return v.([]float64)
}

// GetColumnN from pool, with 0 value
func GetColumnN(n int) []float64 {
	v := colPool.Get()
	if v == nil {
		return make([]float64, n)
	}
	col := v.([]float64)

	for i := 0; i < n; i++ {
		col = append(col, 0)
	}
	return col
}

// PutColumn to pool
func PutColumn(r []float64) {
	colPool.Put(r[:0])
}

// ---------- int64 ----------

// GetInt64 from pool
func GetInt64() []int64 {
	v := int64sPool.Get()
	if v == nil {
		return make([]int64, 0, 100)
	}
	return v.([]int64)
}

// GetInt64N from pool, with 0 value
func GetInt64N(n int) []int64 {
	v := int64sPool.Get()
	if v == nil {
		return make([]int64, n)
	}
	arr := (v.([]int64))

	for i := 0; i < n; i++ {
		arr = append(arr, 0)
	}
	return arr
}

// PutInt64 to pool
func PutInt64(r []int64) {
	int64sPool.Put(r[:0])
}

// ---------- int ----------

// GetInt from pool
func GetInt() []int {
	v := intsPool.Get()
	if v == nil {
		return make([]int, 0, 100)
	}
	return v.([]int)
}

// GetIntN from pool, with 0 value
func GetIntN(n int) []int {
	v := intsPool.Get()
	if v == nil {
		return make([]int, n)
	}
	arr := v.([]int)

	for i := 0; i < n; i++ {
		arr = append(arr, 0)
	}
	return arr
}

// PutInt to pool
func PutInt(r []int) {
	intsPool.Put(r[:0])
}

// ---------- strings ----------

// GetStrings from pool
func GetStrings() []string {
	v := intsPool.Get()
	if v == nil {
		return make([]string, 0, 100)
	}
	return v.([]string)
}

// GetStringsN from pool, with 0 value
func GetStringsN(n int) []string {
	v := strsPool.Get()
	if v == nil {
		return make([]string, n)
	}
	arr := v.([]string)

	for i := 0; i < n; i++ {
		arr = append(arr, "")
	}
	return arr
}

// PutStrings to pool
func PutStrings(r []string) {
	strsPool.Put(r[:0])
}
