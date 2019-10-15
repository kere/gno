package util

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	b := []byte("abc")
	v := string(CRC64Token(b))
	if v != "xJqXpEgAHQ3" {
		t.Fatal(v)
	}
	v = string(CRC32Token(b))
	if v != "mIVkY" {
		t.Fatal(v)
	}
	v = UUID() //1301265771 4217449575
	fmt.Println(v)
}

func TestIntZip(t *testing.T) {
	score := int64(11706300000)
	num := uint64(score)
	v := IntZipTo62(num)
	str := string(v)
	if str != "koreMc" {
		t.Fatal("IntZipTo62", str)
	}
	n := UnIntZip(str, BaseChars)
	if uint64(n) != num {
		t.Fatal(str, n)
	}

	table := make([]byte, len(BaseChars))
	copy(table, BaseChars)
	table = append(table, []byte("-=[];',./~!@#$%^&*()_+{}:\"<>?")...)
	v = IntZipTo(num, table)
	str = string(v)
	if str != "vUD[*1" {
		t.Fatal("IntZipBase", str)
	}
	n = UnIntZip(str, table)
	if uint64(n) != num {
		t.Fatal(str, n, string(table))
	}

	num = 3323724002
	v = IntZipTo62(num)
	str = string(v)
	n = UnIntZip(str, BaseChars)
	if uint64(n) != num {
		t.Fatal(str, n)
	}

	num = 1
	v = IntZipTo62(num)
	str = string(v)
	n = UnIntZip(str, BaseChars)
	if uint64(n) != num {
		t.Fatal(str, n)
	}
	num = 62
	v = IntZipTo62(num)
	str = string(v)
	n = UnIntZip(str, BaseChars)
	if uint64(n) != num {
		t.Fatal(str, n)
	}
}
