package util

import (
	"sort"
)

// Int64s attaches the methods of Interface to []int64, sorting in increasing order.
type Int64s []int64

func (p Int64s) Len() int           { return len(p) }
func (p Int64s) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64s) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Int64s) Sort() { sort.Sort(p) }

// Search searches for x in a sorted slice of ints and returns the index
func (p Int64s) Search(v int64) int {
	return sort.Search(len(p), func(i int) bool { return p[i] >= v })
}

// IndexOf find index
func (p Int64s) IndexOf(v int64) int {
	return IndexOfInt64s(p, v)
}

// IndexOfInt64s 排序法
func IndexOfInt64s(arr []int64, v int64) int {
	l := len(arr)
	isdesc := false
	if l > 2 {
		if arr[0] > arr[1] {
			isdesc = true
		}
	}
	i, isok := getSortedInt64I(v, arr, 0, len(arr)-1, isdesc)
	if isok {
		return i
	}
	return -1
}

// SearchInt64s 排序法
func SearchInt64s(arr []int64, v int64) (int, bool) {
	l := len(arr)
	isdesc := false
	if l > 2 {
		if arr[0] > arr[1] {
			isdesc = true
		}
	}
	return getSortedInt64I(v, arr, 0, len(arr)-1, isdesc)
}

func getSortedInt64I(val int64, arr []int64, b, e int, isdesc bool) (int, bool) {
	//超出边界
	if val < arr[b] {
		return b, false
	} else if val > arr[e] {
		return e + 1, false // 以插入位置为准，所以+1
	}

	switch {
	case arr[b] == val:
		return b, true
	case arr[e] == val:
		return e, true
	}
	diff := e - b
	if diff == 0 {
		return e, false
	} else if diff == 1 {
		return e, false
	} else if diff < 0 {
		return b, false
	}

	l := diff + 1
	i := b + l/2
	v := arr[i]

	if v == val {
		return i, true
	} else if val < v {
		if isdesc {
			return getSortedInt64I(val, arr, i+1, e, isdesc)
		}
		// small zone
		return getSortedInt64I(val, arr, b, i-1, isdesc)
	}
	// v < val
	if isdesc {
		return getSortedInt64I(val, arr, b, i-1, isdesc)
	}
	return getSortedInt64I(val, arr, i+1, e, isdesc)
}

// IndexOfInts 排序法
func IndexOfInts(arr []int, v int) int {
	l := len(arr)
	isdesc := false
	if l > 2 {
		if arr[0] > arr[1] {
			isdesc = true
		}
	}
	i, isok := getSortedIntI(v, arr, 0, len(arr)-1, isdesc)
	if isok {
		return i
	}
	return -1
}

// SearchInts 排序法
func SearchInts(arr []int, v int) (int, bool) {
	l := len(arr)
	isdesc := false
	if l > 2 {
		if arr[0] > arr[1] {
			isdesc = true
		}
	}
	return getSortedIntI(v, arr, 0, len(arr)-1, isdesc)
}

func getSortedIntI(val int, arr []int, b, e int, isdesc bool) (int, bool) {
	//超出边界
	if val < arr[b] {
		return b, false
	} else if val > arr[e] {
		return e + 1, false // 以插入位置为准，所以+1
	}

	switch {
	case arr[b] == val:
		return b, true
	case arr[e] == val:
		return e, true
	}
	diff := e - b
	if diff == 0 {
		return e, false
	} else if diff == 1 {
		return e, false
	} else if diff < 0 {
		return b, false
	}

	l := diff + 1
	i := b + l/2
	v := arr[i]

	if v == val {
		return i, true
	} else if val < v {
		if isdesc {
			return getSortedIntI(val, arr, i+1, e, isdesc)
		}
		// small zone
		return getSortedIntI(val, arr, b, i-1, isdesc)
	}
	// v < val
	if isdesc {
		return getSortedIntI(val, arr, b, i-1, isdesc)
	}
	return getSortedIntI(val, arr, i+1, e, isdesc)
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
