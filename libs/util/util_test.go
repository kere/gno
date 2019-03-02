package util

import (
	"fmt"
	"testing"
)

func Test_Func(t *testing.T) {
	s := append(BaseChars, []byte("-=[];',./'~!@#$%^&*()_+{}:\"<>?")...)
	score := int64(11706300000)
	num := uint64(score)
	v := IntZipTo62(num)
	str := string(v)
	if str != "cMerok" {
		t.Fatal("IntZipTo64", str)
	}

	v = IntZipBaseStr(num, s)
	str = string(v)
	if str != "1'BxHs" {
		t.Fatal("IntZipBase", str)
	}

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

func Test_MapData(t *testing.T) {
	mapData := MapData{}
	mapData["json"] = `[{"p":3},{"p":5},{"p":10}]`

	var v []map[string]int
	err := mapData.JSONParse("json", &v)
	if err != nil {
		t.Fatal(err)
	}

	var v1 []MapData
	err = mapData.JSONParse("json", &v1)
	if err != nil {
		t.Fatal(err)
	}

	mapData = MapData{}
	mapData["a1"] = 1
	mapData["a2"] = 1

	mapData.ArgPlus("a2", 1)
	if mapData.Float("a2") != 2 {
		t.Fatal("ArgPlus failed")
	}

	mapDataCloned := mapData.Clone()
	mapDataCloned.ArgPlus("a2", 1)

	if mapDataCloned.Int("a2")-mapData.Int("a2") != 1 {
		t.Fatal("clone failed")
	}
}

func TestSort(t *testing.T) {
	arr := []int64{10, 322, 3, 43, 65, 30, 230, 44, 56, 76, 20, 430, 659}
	arrSort := Int64s(arr)
	arrSort.Sort()

	index := IndexOfInt64s(arr, 43)
	if index != 4 {
		t.Fatal("index:", index)
	}
	index = IndexOfInt64s(arr, 45)
	if index != -1 {
		t.Fatal("index:", index)
	}

	fmt.Println()
	var isok bool
	arr = []int64{3, 10, 20, 30, 43, 44, 45, 56, 65, 76, 230, 322, 430, 659}
	fmt.Println(arr)
	index, isok = SearchInt64s(arr, 50)
	// if index != 7 || isok {
	// 	t.Fatal(index)
	// }
	fmt.Println("index:", index)
	fmt.Println("---- 50")

	index, isok = SearchInt64s(arr, 57)
	fmt.Println("index:", index)
	fmt.Println("---- 57")

	fmt.Println()
	arr = []int64{3, 10, 20, 30, 43, 44, 45, 49, 56, 65, 76, 230, 322, 430, 659}
	fmt.Println(arr)
	index, isok = SearchInt64s(arr, 50)
	if index != 8 || isok {
		t.Fatal(index)
	}
	fmt.Println("---- 50")
	index, isok = SearchInt64s(arr, 48)
	if index != 7 || isok {
		t.Fatal(index)
	}
	fmt.Println("---- 48")
}
