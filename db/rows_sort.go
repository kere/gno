package db

import "sort"

// RowsSorted class
type RowsSorted struct {
	maprows MapRows
	Field   string
	Method  string
}

// NewRowsSorted func
func NewRowsSorted(maprows MapRows, field string) RowsSorted {
	return RowsSorted{maprows: maprows, Field: field, Method: "int"}
}

// Sort sort.Interface.
func (s *RowsSorted) Sort() {
	sort.Sort(s)
}

// Reverse sort.Interface.
func (s *RowsSorted) Reverse() {
	sort.Sort(sort.Reverse(s))
}

// Len is part of sort.Interface.
func (s *RowsSorted) Len() int {
	return len(s.maprows)
}

// Swap is part of sort.Interface.
func (s *RowsSorted) Swap(i, j int) {
	s.maprows[i], s.maprows[j] = s.maprows[j], s.maprows[i]
}

// Less func.
func (s *RowsSorted) Less(i, j int) bool {
	if s.Method == "int" {
		return s.maprows[i].Int(s.Field) < s.maprows[j].Int(s.Field)
	}
	return s.maprows[i].Float(s.Field) < s.maprows[j].Float(s.Field)
}

// IndexOfInt func.
func (s *RowsSorted) IndexOfInt(v int) int {
	isdesc := false
	if s.maprows.Len() > 2 {
		if s.maprows[0].Int(s.Field) > s.maprows[1].Int(s.Field) {
			isdesc = true
		}
	}
	i, isok := getSortedI(s.Field, v, s.maprows, 0, s.maprows.Len()-1, isdesc)
	if isok {
		return i
	}
	return -1
}

// FindOfInt func.
func (s *RowsSorted) FindOfInt(v int) (MapRow, bool) {
	i := s.IndexOfInt(v)
	if i < 0 {
		return nil, false
	}
	return s.maprows[i], true
}

// IndexOfFloat func.
func (s *RowsSorted) IndexOfFloat(v float64) int {
	isdesc := false
	if s.maprows.Len() > 2 {
		if s.maprows[0].Float(s.Field) > s.maprows[1].Float(s.Field) {
			isdesc = true
		}
	}
	i, isok := getSortedFloatI(s.Field, v, s.maprows, 0, s.maprows.Len()-1, isdesc)
	if isok {
		return i
	}
	return -1
}

// FindOfFloat func.
func (s *RowsSorted) FindOfFloat(v float64) (MapRow, bool) {
	i := s.IndexOfFloat(v)
	if i < 0 {
		return nil, false
	}
	return s.maprows[i], true
}

func getSortedFloatI(field string, val float64, arr MapRows, b, e int, isdesc bool) (int, bool) {
	bVal := arr[b].Float(field)
	eVal := arr[e].Float(field)
	//超出边界
	if val < bVal {
		return b, false
	} else if val > eVal {
		return e + 1, false // 以插入位置为准，所以+1
	}

	switch {
	case bVal == val:
		return b, true
	case eVal == val:
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
	v := arr[i].Float(field)

	if v == val {
		return i, true
	} else if val < v {
		if isdesc {
			return getSortedFloatI(field, val, arr, i+1, e, isdesc)
		}
		// small zone
		return getSortedFloatI(field, val, arr, b, i-1, isdesc)
	}
	// v < val
	if isdesc {
		return getSortedFloatI(field, val, arr, b, i-1, isdesc)
	}
	return getSortedFloatI(field, val, arr, i+1, e, isdesc)
}

func getSortedI(field string, val int, arr MapRows, b, e int, isdesc bool) (int, bool) {
	bVal := arr[b].Int(field)
	eVal := arr[e].Int(field)
	//超出边界
	if val < bVal {
		return b, false
	} else if val > eVal {
		return e + 1, false // 以插入位置为准，所以+1
	}

	switch {
	case bVal == val:
		return b, true
	case eVal == val:
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
	v := arr[i].Int(field)

	if v == val {
		return i, true
	} else if val < v {
		if isdesc {
			return getSortedI(field, val, arr, i+1, e, isdesc)
		}
		// small zone
		return getSortedI(field, val, arr, b, i-1, isdesc)
	}
	// v < val
	if isdesc {
		return getSortedI(field, val, arr, b, i-1, isdesc)
	}
	return getSortedI(field, val, arr, i+1, e, isdesc)
}
