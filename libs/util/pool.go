package util

import "sync"

var (
	rowPool sync.Pool
	colPool sync.Pool
)

// GetRow from pool
func GetRow() []float64 {
	v := rowPool.Get()
	return (v.([]float64))[:0]
}

// GetRowN from pool
func GetRowN(n int) []float64 {
	v := colPool.Get()
	if v == nil {
		return make([]float64, n)
	}
	row := v.([]float64)
	l := len(row)
	if n < l {
		return row[:n+1]
	}

	for i := l - 1; i < n; i++ {
		row = append(row, 0)
	}
	return row
}

// GetRowN0 from pool, with 0 value
func GetRowN0(n int) []float64 {
	v := colPool.Get()
	if v == nil {
		return make([]float64, n)
	}
	row := (v.([]float64))[:0]

	for i := 0; i < n; i++ {
		row = append(row, 0)
	}
	return row
}

// PutRow to pool
func PutRow(r []float64) {
	rowPool.Put(r)
}

// GetColumn from pool
func GetColumn() []float64 {
	v := colPool.Get()
	if v == nil {
		return []float64{}
	}
	return (v.([]float64))[:0]
}

// GetColumnN from pool
func GetColumnN(n int) []float64 {
	v := colPool.Get()
	if v == nil {
		return make([]float64, n)
	}
	col := v.([]float64)
	l := len(col)
	if n < l {
		return col[:n+1]
	}

	for i := l - 1; i < n; i++ {
		col = append(col, 0)
	}
	return col
}

// GetColumnN0 from pool, with 0 value
func GetColumnN0(n int) []float64 {
	v := colPool.Get()
	if v == nil {
		return make([]float64, n)
	}
	col := (v.([]float64))[:0]

	for i := 0; i < n; i++ {
		col = append(col, 0)
	}
	return col
}

// PutColumn to pool
func PutColumn(r []float64) {
	colPool.Put(r)
}
