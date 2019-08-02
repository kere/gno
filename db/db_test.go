package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/conf"
	"github.com/lib/pq"
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

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestInsert(t *testing.T) {
	Current().Exec("truncate table " + table)
	ins := NewInsert(table)
	row := MapRow{"a": "TestA", "b": -1, "created_at": time.Now(), "vals": []float64{1, 2}, "ints": []int{1, 2}, "v_json": nil}
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
	if row.String("a") != "TestA" {
		t.Fatal(row)
	}
}

func TestInsertM(t *testing.T) {
	ins := NewInsert(table)
	now := time.Now()
	l := 7
	rows := NewMapRows(l)
	for i := 0; i < l; i++ {
		rows[i] = MapRow{"a": fmt.Sprint("a", i), "b": i, "created_at": now, "vals": []float64{float64(i) + 0.1, float64(i)}, "ints": []int{1 + i, 2 + i}, "v_json": User{"tom", 22 + i}, "strings": []string{fmt.Sprint("s", i), fmt.Sprint("s", i+1)}}
	}

	ins = NewInsert(table)
	r, err := ins.InsertM(rows)
	n, _ := r.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if n != 7 {
		t.Fatal("effected rows:", n)
	}

	q := NewQuery(table)
	dat, err := q.Query()
	if err != nil {
		t.Fatal(err)
	}

	PrintDataSet(&dat)
	_, err = dat.TimeAt(1, "created_at")
	if err != nil {
		t.Fatal(err)
	}

	src, _ := dat.BytesAt(0, "vals")

	arr := pq.Float64Array{}
	err = arr.Scan(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(arr) != 2 || arr[1] != 2 {
		t.Fatal()
	}
}

func TestDBQueryRows(t *testing.T) {
	q := NewQuery(table)
	rows, err := q.QueryRows()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) < 4 {
		t.Fatal(len(rows))
	}

	row := rows[3]
	str := row.String("a")
	if str != "a2" {
		t.Fatal(str)
	}

	strs := row.Strings("strings")
	if len(strs) != 2 || strs[1] != "s3" {
		t.Fatal(strs)
	}

	ints := row.Ints("ints")
	if len(ints) != 2 || ints[1] != 4 {
		t.Fatal(ints)
	}

	floats := row.Floats("vals")
	if len(floats) != 2 || floats[1] != 2 {
		t.Fatal(floats)
	}
}

func TestDBQuery(t *testing.T) {
	q := NewQuery(table)
	dat, err := q.Query()
	if err != nil {
		t.Fatal(err)
	}
	if dat.Len() < 4 {
		t.Fatal(dat.Len())
	}

	i := 3
	str, _ := dat.StrAt(i, "a")
	if str != "a2" {
		t.Fatal(str)
	}

	strs, _ := dat.StrsAt(i, "strings")
	if len(strs) != 2 || strs[1] != "s3" {
		t.Fatal(strs)
	}

	ints, _ := dat.IntsAt(i, "ints")
	if len(ints) != 2 || ints[1] != 4 {
		t.Fatal(ints)
	}
	int64s, _ := dat.Int64sAt(i, "ints")
	if len(int64s) != 2 || int64s[1] != 4 {
		t.Fatal(int64s)
	}

	floats, _ := dat.FloatsAt(i, "vals")
	if len(floats) != 2 || floats[1] != 2 {
		t.Fatal(floats)
	}
}
