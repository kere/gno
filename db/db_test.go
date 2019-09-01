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

type DataRow struct {
	A         string    `json:"a"`
	B         int       `json:"b"`
	Vals      []float64 `json:"vals"`
	Ints      []int     `json:"ints"`
	VJSON     User      `json:"v_json"`
	CreatedAt time.Time `json:"created_at"`
}

func TestConvert(t *testing.T) {
	row := MapRow{"a": "TestA", "b": -1, "vals": []float64{1, 2}, "ints": []int{1, 2}, "v_json": "{\"name\":\"tome\", \"age\": 22}", "created_at": time.Now()}
	vo := DataRow{}

	Row2VO(row, &vo)

	if vo.A != row.String("a") {
		t.Fatal()
	}
	if vo.B != row.Int("b") {
		t.Fatal()
	}
	vals := row.Floats("vals")
	if len(vo.Vals) != len(vals) || vo.Vals[0] != vals[0] {
		t.Fatal()
	}

	if vo.VJSON.Name != "tome" || vo.VJSON.Age != 22 {
		t.Fatal()
	}
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

func TestQueryCache(t *testing.T) {
	SetCache(cache.CurrentCache())
	q := NewQuery(table).Cache()
	q.Query()

	key := querybuildCacheKey(q, 0)
	if isok, _ := cacheIns.IsExists(key); !isok {
		t.Fatal("cache not found")
	}

	dat, _ := q.Query()
	if dat.Len() < 4 {
		t.Fatal(dat.Len())
	}
	i := 3
	str, _ := dat.StrAt(i, "a")
	if str != "a2" {
		t.Fatal(str)
	}
}

type testUsr struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type testVO struct {
	A       string    `json:"a"`
	B       int       `json:"b"`
	Vals    []float64 `json:"vals"`
	Strs    []string  `json:"strings"`
	Ints    []int     `json:"ints"`
	Dat     testUsr   `json:"v_json"`
	Created time.Time `json:"created_at"`
}

func (v *testVO) Table() string {
	return "temp"
}

func TestIVO(t *testing.T) {
	now := time.Now().AddDate(0, -10, 0)
	vo := testVO{A: "vo01", B: 100, Vals: []float64{11.1, 12.2, 13.3}, Strs: []string{"vo1", "vo2"}, Ints: []int{10, 11, 12}, Dat: testUsr{"Ins", 10}, Created: now}

	row := VO2InsertMapRow(&vo)
	if len(row) == 0 || len((row["vals"].([]float64))) != 3 {
		t.Fatal(row)
	}

	err := VOCreate(&vo)
	if err != nil {
		t.Fatal(err)
	}

	q := NewQuery(vo.Table()).Where("a=? and b=?", vo.A, vo.B)
	dat, err := q.QueryOne()
	if err != nil {
		t.Fatal(err)
	}
	if dat.IsEmpty() {
		t.Fatal("failed to insert")
	}
	if dat.String("a") != vo.A || dat.Int("b") != vo.B {
		t.Fatal(dat)
	}

	u := NewUpdate(vo.Table()).Where("a=? and b=?", vo.A, vo.B)

	vo = testVO{A: "vo02", B: 300, Vals: []float64{31.1, 32.2, 33.3}, Strs: []string{"vo30", "vo31"}, Ints: []int{30, 31, 32}, Dat: testUsr{"Upd", 30}, Created: now.AddDate(1, 0, 0)}
	row = VO2InsertMapRow(&vo)
	_, err = u.Update(row)
	if err != nil {
		t.Fatal(err)
	}

}
