package db

import (
	"fmt"
	"sort"
	"testing"
)

func TestDataRow(t *testing.T) {
	row1 := MapRow{
		"name": "tom",
		"age":  22,
		"amt":  32.23,
		"arr":  []string{"a", "b", "c"},
	}
	row2 := MapRow{
		"name": "tom",
		"age":  22,
		"amt":  32.23,
		"arr":  []string{"a", "b", "c"},
	}

	dat := row1.ChangedData(row2)
	if len(dat) != 0 {
		t.Fatal(dat)
	}

	row2 = MapRow{
		"name": "tom1",
		"age":  21,
		"amt":  32.239,
		"arr":  []string{"a", "b", "c"},
	}
	dat = row1.ChangedData(row2)
	if len(dat) != 3 || dat.String("name") != "tom1" || dat.Float("amt") != row2.Float("amt") {
		t.Fatal(dat)
	}

	row2 = MapRow{
		"name": "tom",
		"age":  22,
		"amt":  32.23,
		"arr":  []string{"e", "b", "c"},
	}
	dat = row1.ChangedData(row2)
	fmt.Println(dat)
}

func TestDataSet(t *testing.T) {
	var row MapRow
	arr := MapRows{}
	row = MapRow{"id": 0, "name": "tom1", "age": 20}
	arr = append(arr, row)
	row = MapRow{"id": 1, "name": "tom1", "age": 21}
	arr = append(arr, row)
	row = MapRow{"id": 2, "name": "tom2", "age": 22}
	arr = append(arr, row)
	row = MapRow{"id": 3, "name": "tom3", "age": 23}
	arr = append(arr, row)
	row = MapRow{"id": 4, "name": "tom4", "age": 24}
	arr = append(arr, row)
	row = MapRow{"id": 5, "name": "tom5", "age": 25}
	arr = append(arr, row)

	s := NewRowsSorted(arr, "id")

	if !sort.IsSorted(&s) {
		t.Fatal()
	}
	i := s.IndexOfInt(3)
	if i != 3 {
		t.Fatal(i)
	}
}
