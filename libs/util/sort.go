package util

import (
	"sort"
)

// Int64sOrder attaches the methods of Interface to []int64, sorting in increasing order.
type Int64sOrder []int64

func (p Int64sOrder) Len() int           { return len(p) }
func (p Int64sOrder) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64sOrder) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Int64sOrder) Sort() { sort.Sort(p) }

// Search searches for x in a sorted slice of ints and returns the index
func (p Int64sOrder) Search(v int64) int {
	return sort.Search(len(p), func(i int) bool { return p[i] >= v })
}

// IndexOf find index, if not found, return -1
func (p Int64sOrder) IndexOf(v int64) int {
	i := p.Search(v)
	if p[i] != v {
		return -1
	}
	return i
}

// Int64Indexs db int64 orded column
type Int64Indexs [][2]int64 // 排序列[2]int64, 0: val, 1: orgin index

// Len is part of sort.Interface.
func (s Int64Indexs) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s Int64Indexs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less sort.
func (s Int64Indexs) Less(i, j int) bool {
	return s[i][0] < s[j][0]
}

// IndexOf not found return -1
func (s Int64Indexs) IndexOf(v int64) int {
	index := sort.Search(len(s), func(i int) bool {
		return s[i][0] >= v
	})
	if s[index][0] == v {
		return index
	}
	return -1
}

// Search data
func (s Int64Indexs) Search(v int64) int {
	return sort.Search(len(s), func(i int) bool {
		return s[i][0] >= v
	})
}

// IndexOfStrs f
func IndexOfStrs(arr []string, v string) int {
	n := len(arr)
	for i := 0; i < n; i++ {
		if arr[i] == v {
			return i
		}
	}
	return -1
}

// Float32Slice attaches the methods of Interface to []float32, sorting in increasing order
// (not-a-number values are treated as less than other values).
type Float32Slice []float32

func (p Float32Slice) Len() int           { return len(p) }
func (p Float32Slice) Less(i, j int) bool { return p[i] < p[j] || isNaN32(p[i]) && !isNaN32(p[j]) }
func (p Float32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// isNaN32 is a copy of math.IsNaN to avoid a dependency on the math package.
func isNaN32(f float32) bool {
	return f != f
}

// Sort is a convenience method.
func (p Float32Slice) Sort() { sort.Sort(p) }

// Float32s sorts a slice of float32s in increasing order
// (not-a-number values are treated as less than other values).
func Float32s(a []float32) { sort.Sort(Float32Slice(a)) }
