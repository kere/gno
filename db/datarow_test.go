package db

import (
	"fmt"
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
}