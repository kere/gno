package db

import (
	"fmt"
	"sort"
	"testing"
)

func TestDataRow(t *testing.T) {
	row1 := DataRow{
		"name": "tom",
		"age":  22,
		"amt":  32.23,
		"arr":  []string{"a", "b", "c"},
	}
	row2 := DataRow{
		"name": "tom",
		"age":  22,
		"amt":  32.23,
		"arr":  []string{"a", "b", "c"},
	}

	dat := row1.ChangedData(row2)
	if len(dat) != 0 {
		t.Fatal(dat)
	}

	row2 = DataRow{
		"name": "tom1",
		"age":  21,
		"amt":  32.239,
		"arr":  []string{"a", "b", "c"},
	}
	dat = row1.ChangedData(row2)
	if len(dat) != 3 || dat.String("name") != "tom1" || dat.Float("amt") != row2.Float("amt") {
		t.Fatal(dat)
	}

	row2 = DataRow{
		"name": "tom",
		"age":  22,
		"amt":  32.23,
		"arr":  []string{"e", "b", "c"},
	}
	dat = row1.ChangedData(row2)
	fmt.Println(dat)

	row2["uint"] = uint(126)
	if row2.Int("uint") != 126 {
		t.Fatal("uint failed")
	}
}

func TestDataSet(t *testing.T) {
	var row DataRow
	arr := DataSet{}
	row = DataRow{"id": 0, "name": "tom1", "age": 20}
	arr = append(arr, row)
	row = DataRow{"id": 1, "name": "tom1", "age": 21}
	arr = append(arr, row)
	row = DataRow{"id": 2, "name": "tom2", "age": 22}
	arr = append(arr, row)
	row = DataRow{"id": 3, "name": "tom3", "age": 23}
	arr = append(arr, row)
	row = DataRow{"id": 4, "name": "tom4", "age": 24}
	arr = append(arr, row)
	row = DataRow{"id": 5, "name": "tom5", "age": 25}
	arr = append(arr, row)

	s := NewDataSetSorted(arr, "id")

	if !sort.IsSorted(&s) {
		t.Fatal()
	}
	i := s.IndexOfInt(3)
	if i != 3 {
		t.Fatal(i)
	}
}
