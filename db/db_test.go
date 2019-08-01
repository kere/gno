package db

import (
	"fmt"
	"testing"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/conf"
)

var (
	table = "temp"
)

func init() {
	file := "../httpd/example/app/app.conf"
	c := conf.Load(file)
	Init("db", c.GetConf("db"))
	cache.Init(c.GetConf("cache"))
}

func TestInsert(t *testing.T) {
	Current().Exec("truncate table " + table)
	ins := NewInsert(table)
	row := MapRow{"a": "TestA", "b": -1, "c": "TestC"}
	r, err := ins.Insert(row)
	n, _ := r.RowsAffected()
	if err != nil || n != 1 {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	q := NewQuery(table)
	rows, err := q.QueryRows()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rows)

	dataset, err := q.Query()
	if err != nil {
		t.Fatal(err)
	}
	row := dataset.MapRowAt(0)
	if row.String("c") != "TestC" {
		t.Fatal(row)
	}
}

func TestInsertM(t *testing.T) {
	ins := NewInsert(table)
	l := 7
	rows := make([]MapRow, l)
	for i := 0; i < l; i++ {
		rows[i] = MapRow{"a": fmt.Sprint("a", i), "b": i, "c": fmt.Sprint("c", i)}
	}

	ins = NewInsert(table)
	r, err := ins.InsertM(rows)
	n, _ := r.RowsAffected()
	if err != nil || n != 7 {
		t.Fatal(err)
	}
}
