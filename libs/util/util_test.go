package util

import (
	"testing"
)

func Test_Func(t *testing.T) {
	s := append(BaseChars, []byte("-=[];',./'~!@#$%^&*()_+{}:\"<>?")...)
	score := 11706300000
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
