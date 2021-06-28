package util

import (
	"sync"
)

type poolClas struct {
	sync.Pool
	lock sync.Mutex
}

var (
	itemsPool  = poolClas{}
	floatsPool = poolClas{}
	intsPool   = poolClas{}
	strsPool   = poolClas{}
	int64sPool = poolClas{}
	bytesPool  = poolClas{}
)

func ParsePoolArgs(args []int, minCap int) (l int, capN int) {
	n := len(args)
	switch {
	case n == 0:
		l, capN = 0, 0
	case n == 1:
		l, capN = args[0], args[0]
	default:
		l, capN = args[0], args[1]
	}
	if capN < minCap {
		capN = minCap
	}
	return
}

// GetBytes for bit setup
func GetBytes(args ...int) []byte {
	l, capN := ParsePoolArgs(args, 50)
	bytesPool.lock.Lock()
	defer bytesPool.lock.Unlock()
	v := bytesPool.Get()
	if v == nil {
		return make([]byte, l, capN)
	}
	row := v.([]byte)

	for i := 0; i < l; i++ {
		row = append(row, 0)
	}
	return row
}

// PutBytes for bit setup
func PutBytes(arr []byte) {
	if cap(arr) == 0 {
		return
	}
	bytesPool.Put(arr[:0])
}

// GetFloats from pool
func GetFloats(args ...int) []float64 {
	l, capN := ParsePoolArgs(args, 50)
	floatsPool.lock.Lock()
	defer floatsPool.lock.Unlock()

	v := floatsPool.Get()

	if v == nil {
		return make([]float64, l, capN)
	}
	row := v.([]float64)

	for i := 0; i < l; i++ {
		row = append(row, 0)
	}
	return row
}

// PutfloatsPool to pool
func PutFloats(r []float64) {
	if cap(r) == 0 {
		return
	}
	floatsPool.Put(r[:0])
}

// ---------- int64 ----------

// GetInt64s from pool
func GetInt64s(args ...int) []int64 {
	l, capN := ParsePoolArgs(args, 50)
	int64sPool.lock.Lock()
	defer int64sPool.lock.Unlock()

	v := int64sPool.Get()
	if v == nil {
		return make([]int64, l, capN)
	}

	arr := (v.([]int64))
	for i := 0; i < l; i++ {
		arr = append(arr, 0)
	}
	return arr
}

// PutInt64s to pool
func PutInt64s(ints []int64) {
	if cap(ints) == 0 {
		return
	}
	int64sPool.Put(ints[:0])
}

// ---------- int ----------

// GetInts from pool, with 0 value
func GetInts(args ...int) []int {
	l, capN := ParsePoolArgs(args, 50)
	intsPool.lock.Lock()
	defer intsPool.lock.Unlock()
	v := intsPool.Get()
	if v == nil {
		return make([]int, l, capN)
	}
	arr := v.([]int)

	for i := 0; i < l; i++ {
		arr = append(arr, 0)
	}
	return arr
}

// PutInt to pool
func PutInts(ints []int) {
	if cap(ints) == 0 {
		return
	}
	intsPool.Put(ints[:0])
}

// ---------- strings ----------

// GetStrings from pool
func GetStrings(args ...int) []string {
	l, capN := ParsePoolArgs(args, 50)

	strsPool.lock.Lock()
	defer strsPool.lock.Unlock()

	v := strsPool.Get()
	if v == nil {
		return make([]string, l, capN)
	}
	arr := v.([]string)

	for i := 0; i < l; i++ {
		arr = append(arr, "")
	}
	return arr
}

// PutStrings to pool
func PutStrings(strs []string) {
	if cap(strs) == 0 {
		return
	}
	strsPool.Put(strs[:0])
}

// GetItems from pool
func GetItems(args ...int) []interface{} {
	l, capN := ParsePoolArgs(args, 50)
	itemsPool.lock.Lock()
	defer itemsPool.lock.Unlock()

	v := itemsPool.Get()

	if v == nil {
		return make([]interface{}, l, capN)
	}
	row := v.([]interface{})

	for i := 0; i < l; i++ {
		row = append(row, nil)
	}
	return row
}

// PutItems to pool
func PutItems(r []interface{}) {
	if cap(r) == 0 {
		return
	}
	itemsPool.Put(r[:0])
}
