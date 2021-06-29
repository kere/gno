package db

import (
	"testing"
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
	PutRow(row)
	dat, _ = q.Query()
	if dat.Len() != 1 || dat.Columns[0][0].(string) != row[0].(string) {
		PrintDataSet(&dat)
		t.Fatal(row)
	}
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
