package util

import (
	"fmt"
	"testing"
)

func Test_Func(t *testing.T) {
	num2 := 1.712774821
	val := Round(num2, 2)
	if val != 1.71 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 3)
	if val != 1.713 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 0)
	if val != 2.0 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 4)
	if val != 1.7128 {
		t.Fatal("round failed", val)
	}
	num2 = -1.712774821
	val = Round(num2, 2)
	if val != -1.71 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 3)
	if val != -1.713 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 0)
	if val != -2.0 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 4)
	if val != -1.7128 {
		t.Fatal("round failed", val)
	}
}

func TestSort(t *testing.T) {
	arr := []int64{10, 322, 3, 43, 65, 30, 230, 44, 56, 76, 20, 430, 659}
	arrSort := Int64sOrder(arr)
	arrSort.Sort()

	index := arrSort.IndexOf(43)
	if index != 4 {
		t.Fatal("index:", index)
	}
	index = arrSort.IndexOf(45)
	if index != -1 {
		t.Fatal("index:", index)
	}

	arr = []int64{3, 10, 20, 30, 43, 44, 45, 56, 65, 76, 230, 322, 430, 659}
	arrSort = Int64sOrder(arr)
	index = arrSort.Search(50)

	if index != 7 {
		fmt.Println(arr)
		t.Fatal(index)
	}

	index = arrSort.Search(57)
	if index != 8 {
		fmt.Println(arr)
		t.Fatal(index)
	}

	arr = []int64{3, 10, 20, 30, 43, 44, 45, 49, 56, 65, 76, 230, 322, 430, 659}
	arrSort = Int64sOrder(arr)
	index = arrSort.Search(50)
	if index != 8 {
		fmt.Println(arr)
		t.Fatal(index)
	}
	index = arrSort.Search(48)
	if index != 7 {
		fmt.Println(arr)
		t.Fatal(index)
	}
}

func TestSync(t *testing.T) {
	cpt := NewComputation()
	arr := make([]int, 100)

	counter := 0
	cpt.RunA(100, func(i int) (interface{}, error) {
		arr[i] = i + 1
		return i, nil
	}, func(i int, dat interface{}) {
		counter++
	})
	// fmt.Println(counter, "===========")
	// fmt.Println(arr, "===========")

	if counter != 100 {
		fmt.Println(counter, "===========")
		fmt.Println(arr, "===========")
		t.Fatal(counter)
	}

	s := []byte("110010")
	v := BitStr2Uint(s)
	if v != 50 {
		t.Fatal(v)
	}
	if !IsTrueAtBitUint(v, 1) {
		t.Fatal(s)
	}
	if IsTrueAtBitUint(v, 2) {
		t.Fatal(s)
	}

}
func TestPool(t *testing.T) {
	r := make([]float64, 5)
	for i := 0; i < 5; i++ {
		r[i] = float64(i + 1)
	}

	PutRow(r)
	r = GetRowN(10)
	if len(r) != 10 {
		t.Fatal(r)
	}

	PutRow(r)
	r = GetRowN(10)
	if len(r) != 10 {
		t.Fatal(r)
	}

	PutRow(r)
	r = GetRowN(5)
	if len(r) != 5 {
		t.Fatal(r)
	}
}

func TestCamelCase(t *testing.T) {
	s := "created_at_some"
	str := CamelCase(s)
	if str != "CreatedAtSome" {
		t.Fatal(str)
	}
}
