package db

import "sort"

// DataSetSorted class
type DataSetSorted struct {
	dataset DataSet
	Field   string
	Method  string
}

// NewDataSetSorted func
func NewDataSetSorted(dataset DataSet, field string) DataSetSorted {
	return DataSetSorted{dataset: dataset, Field: field, Method: "int"}
}

// Sort sort.Interface.
func (s *DataSetSorted) Sort() {
	sort.Sort(s)
}

// Reverse sort.Interface.
func (s *DataSetSorted) Reverse() {
	sort.Reverse(s)
}

// Len is part of sort.Interface.
func (s *DataSetSorted) Len() int {
	return len(s.dataset)
}

// Swap is part of sort.Interface.
func (s *DataSetSorted) Swap(i, j int) {
	s.dataset[i], s.dataset[j] = s.dataset[j], s.dataset[i]
}

// Less func.
func (s *DataSetSorted) Less(i, j int) bool {
	if s.Method == "int" {
		return s.dataset[i].Int(s.Field) < s.dataset[j].Int(s.Field)
	}
	return s.dataset[i].Float(s.Field) < s.dataset[j].Float(s.Field)
}

// IndexOfInt func.
func (s *DataSetSorted) IndexOfInt(v int) int {
	isdesc := false
	if s.dataset.Len() > 2 {
		if s.dataset[0].Int(s.Field) > s.dataset[1].Int(s.Field) {
			isdesc = true
		}
	}
	i, isok := getSortedI(s.Field, v, s.dataset, 0, s.dataset.Len()-1, isdesc)
	if isok {
		return i
	}
	return -1
}

// FindOfInt func.
func (s *DataSetSorted) FindOfInt(v int) (DataRow, bool) {
	i := s.IndexOfInt(v)
	if i < 0 {
		return nil, false
	}
	return s.dataset[i], true
}

// IndexOfFloat func.
func (s *DataSetSorted) IndexOfFloat(v float64) int {
	isdesc := false
	if s.dataset.Len() > 2 {
		if s.dataset[0].Float(s.Field) > s.dataset[1].Float(s.Field) {
			isdesc = true
		}
	}
	i, isok := getSortedFloatI(s.Field, v, s.dataset, 0, s.dataset.Len()-1, isdesc)
	if isok {
		return i
	}
	return -1
}

// FindOfFloat func.
func (s *DataSetSorted) FindOfFloat(v float64) (DataRow, bool) {
	i := s.IndexOfFloat(v)
	if i < 0 {
		return nil, false
	}
	return s.dataset[i], true
}

func getSortedFloatI(field string, val float64, arr DataSet, b, e int, isdesc bool) (int, bool) {
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

func getSortedI(field string, val int, arr DataSet, b, e int, isdesc bool) (int, bool) {
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
