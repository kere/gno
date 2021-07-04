package db

import (
	"testing"

	"github.com/kere/gno/libs/util"
	_ "github.com/lib/pq"
)

var (
	table = "test01"
)

func init() {
	config := make(map[string]string)
	config["driver"] = "postgres"
	config["dbname"] = "astock"
	config["user"] = "postgres"
	config["password"] = "123"
	// config["level"] = "all"
	Init("app", config)
	Current().SetLogLevel("all")
}

type User struct {
	Name string
	Age  int
}

func Test_A(t *testing.T) {
	b := NewBuilder(table)
	sql := `create table if not exists test01 (
code                 VARCHAR(20)         not null,
date                 INT4                not null,
a_json                JSONB              null,
values               FLOAT4[]            null
); `
	_, err := b.Exec(sql, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.Exec("truncate table "+table, nil)
	if err != nil {
		t.Fatal(err)
	}

	// 1: 测试读取数据库fields
	q := NewQuery(table)
	dat, err := q.Query()
	if err != nil {
		t.Fatal(err)
	}
	if len(dat.Fields) == 0 || dat.Fields[0] != "code" {
		t.Fatal(dat.Fields)
	}

	// 2: 测试单次插入
	ins := NewInsert(table)
	fields := []string{"code", "date", "a_json", "values"}

	row := GetRow(len(fields))
	row[0] = "code01"
	row[1] = 1
	row[2] = User{Name: "tom", Age: 22}
	row[3] = []float64{1.1, 1.2, 1.3}
	_, err = ins.Insert(fields, row)
	if err != nil {
		t.Fatal(err)
	}
	dat, _ = q.Query()
	if dat.Len() != 1 || dat.Columns[0][0].(string) != row[0].(string) {
		PrintDataSet(&dat)
		t.Fatal(row)
	}
	PutRow(row)
	row = dat.RowAtP(0)

	if string(row[3].([]byte)) != "{1.1,1.2,1.3}" {
		PrintDataSet(&dat)
		t.Fatal(row)
	}

	// 3: 测试多次插入
	// dat = NewDataSet(fields)
	dat = GetDataSet(fields)
	row[0] = "code02"
	row[1] = 2
	row[2] = User{Name: "tom02", Age: 20}
	row[3] = []float64{2.1, 2.2, 2.3}
	dat.AddRow(row)
	row[0] = "code03"
	row[1] = 3
	row[2] = User{Name: "tom03", Age: 32}
	row[3] = []float64{3.1, 3.2, 3.3}
	dat.AddRow(row)
	row[0] = "code04"
	row[1] = 4
	row[2] = User{Name: "tom04", Age: 42}
	row[3] = []float64{4.1, 4.2, 4.3}
	dat.AddRow(row)
	row[0] = "code05"
	row[1] = 2
	row[2] = User{Name: "tom05", Age: 52}
	row[3] = []float64{5.1, 5.2, 5.3}
	dat.AddRow(row)
	_, err = ins.InsertM(&dat)
	if err != nil {
		t.Fatal(err)
	}
	PutDataSet(&dat)
	dat, err = q.Query()
	if err != nil {
		PrintDataSet(&dat)
		t.Fatal(err)
	}
	if dat.Len() != 5 || dat.Columns[0][4].(string) != "code05" || dat.Columns[1][4].(int64) != 2 {
		PrintDataSet(&dat)
		t.Fatal()
	}

	// 4: 测试Update
	u := Current().NewUpdate(table)
	row = GetRow(1)
	row[0] = 5
	result, err := u.Where("code=$1 and date=$2", "code05", 2).Update([]string{"date"}, row)
	if err != nil {
		t.Fatal(err)
	}
	n, _ := result.RowsAffected()
	if n != 1 {
		t.Fatal("update failed")
	}

	dat, _ = q.Query()
	if dat.Len() != 5 || dat.Columns[0][4].(string) != "code05" || dat.Columns[1][4].(int64) != 5 {
		PrintDataSet(&dat)
		t.Fatal()
	}

	// 5: 测试Exists
	e := Current().NewExists(table)
	if e.Where("code=$1 and date=$2", "code05", 5).NotExists() {
		t.Fatal("exists failed")
	}
	// 6: 测试Delete
	del := Current().NewDelete(table)
	r, err := del.Where("code=$1 and date=$2", "code05", 5).Delete()
	n, _ = r.RowsAffected()
	if n != 1 {
		t.Fatal("delete failed")
	}
	if e.Where("code=$1 and date=$2", "code05", 5).Exists() {
		t.Fatal("exists failed")
	}

	// 7: 检查数据
	dat, _ = q.Query()
	if dat.Len() != 4 || dat.Columns[0][3].(string) != "code04" || dat.Columns[1][2].(int64) != 3 {
		PrintDataSet(&dat)
		t.Fatal()
	}
}

func Test_Tx(t *testing.T) {
	tx, err := BeginTx()
	if err != nil {
		t.Fatal(err)
	}

	fields := []string{"code", "date", "a_json", "values"}

	row := GetRow(len(fields))
	row[0] = "code10"
	row[1] = 10
	row[2] = User{Name: "tom", Age: 10}
	row[3] = []float64{10.1, 10.2, 10.3}

	// 1: Insert row commit
	ins := tx.NewInsert(table)
	_, err = ins.Insert(fields, row)
	if err != nil {
		t.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}
	q := Current().NewQuery(table)
	dat, err := q.Query()
	if err != nil {
		t.Fatal(err)
	}
	if dat.Len() != 5 || dat.Columns[0][4].(string) != "code10" || dat.Columns[1][4].(int64) != 10 {
		PrintDataSet(&dat)
		t.Fatal()
	}

	tx, _ = BeginTx()
	ins = tx.NewInsert(table)
	// 2: Insert row rollback
	row[0] = "code11"
	row[1] = 11
	row[2] = User{Name: "tom", Age: 11}
	row[3] = []float64{11.1, 11.2, 11.3}
	_, err = ins.Insert(fields, row)
	if err != nil {
		t.Fatal(err)
	}
	err = tx.Rollback()
	if err != nil {
		t.Fatal(err)
	}

	q = Current().NewQuery(table)
	dat, err = q.Query()
	if err != nil {
		t.Fatal(err)
	}
	if dat.Len() != 5 || dat.Columns[0][4].(string) != "code10" || dat.Columns[1][4].(int64) != 10 {
		PrintDataSet(&dat)
		t.Fatal()
	}
}

func TestDBRow(t *testing.T) {
	dbRow := DBRow{Fields: []string{"code", "date", "ints", "floats", "strings"}}
	dbRow.Values = []interface{}{
		[]byte("a001"),
		int64(1190101),
		[]byte("{1,2,3}"),
		[]byte("{1.1,2.2,3.3}"),
		[]byte(`{a,b,c}`),
	}

	ints, err := dbRow.IntsAt(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(ints) != 3 || ints[2] != 3 {
		t.Fatal(ints)
	}
	floats, err := dbRow.FloatsAtP(3)
	if err != nil {
		t.Fatal(err)
	}
	if len(floats) != 3 || floats[2] != 3.3 {
		t.Fatal(floats)
	}
	util.PutFloats(floats)
	strs, err := dbRow.StringsAt(4)
	if err != nil {
		t.Fatal(err)
	}
	if len(strs) != 3 || strs[2] != "c" {
		t.Fatal(strs)
	}

	dbRow.Values[3] = []byte("{2.2404713e+09,3.8556639e+09,6.096135e+09}")
	ints, err = dbRow.IntsAt(3)
	if err != nil {
		t.Fatal(err)
	}
	if ints[0] != 2240471300 {
		t.Fatal(ints)
	}
	dbRow.Values[3] = []byte("{2.9352145e+10,0,2.9352145e+10}")
	int64s, err := dbRow.Int64sAt(3)
	if err != nil {
		t.Fatal(err)
	}
	if int64s[2] != 2.9352145e+10 || int64s[1] != 0 {
		t.Fatal(int64s)
	}
}
func TestPage(t *testing.T) {
	count := 3
	// 0,1,2
	ds := GetDataSet([]string{"val"}, count)
	for i := 0; i < count; i++ {
		ds.Columns[0][i] = i
	}
	n := ds.EachPage(3, func(page int, ds1 DataSet) bool {
		if page != 1 {
			PrintDataSet(&ds1)
			t.Fatal()
			return false
		}
		return true
	})
	if n != 1 {
		t.Fatal()
	}

	// 0,1,2,3
	row := ds.GetRowP()
	row[0] = 3
	ds.AddRow(row)
	n = ds.EachPage(3, func(page int, ds1 DataSet) bool {
		if page == 1 {
			if ds1.Len() != 3 {
				t.Fatal()
				return false
			}
		} else {
			if ds1.Len() != 1 {
				PrintDataSet(&ds1)
				t.Fatal()
				return false
			}
		}
		return true
	})
	if n != 2 {
		t.Fatal(n)
	}

	// 0,1,2,3,4
	row[0] = 4
	ds.AddRow(row)
	n = ds.EachPage(3, func(page int, ds1 DataSet) bool {
		if page == 2 && ds1.Len() != 2 {
			PrintDataSet(&ds1)
			t.Fatal()
		}
		return true
	})
	if n != 2 {
		t.Fatal(n)
	}

	// 0,1,2,3,4,5
	row[0] = 5
	ds.AddRow(row)
	n = ds.EachPage(3, func(page int, ds1 DataSet) bool {
		if page == 2 && ds1.Len() != 3 {
			PrintDataSet(&ds1)
			t.Fatal()
		}
		return true
	})
	if n != 2 {
		t.Fatal(n)
	}

	// 0,1,2,3,4,5,6
	row[0] = 6
	ds.AddRow(row)
	n = ds.EachPage(3, func(page int, ds1 DataSet) bool {
		if page == 2 && ds1.Len() != 3 {
			PrintDataSet(&ds1)
			t.Fatal()
		}
		if page == 3 && ds1.Len() != 1 {
			PrintDataSet(&ds1)
			t.Fatal()
		}
		return true
	})
	if n != 3 {
		t.Fatal(n)
	}

	PutDataSet(&ds)
	PutRow(row)
}
