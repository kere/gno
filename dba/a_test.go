package dba

import (
	"testing"
)

func init() {
	config := make(map[string]string)
	config["driver"] = "postgres"
	config["dbname"] = "astock"
	config["user"] = "postgres"
	config["password"] = "123"
	Init("app", config)
	Current().SetLogLevel("all")
}

type User struct {
	Name string
	Age  int
}

func Test_A(t *testing.T) {
	b := Builder{}
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
	table := "test01"
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
	_, err = ins.Insert0(fields, row)
	if err != nil {
		t.Fatal(err)
	}
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
	row[1] = 5
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
	if dat.Len() != 5 || dat.Columns[0][4].(string) != "code05" {
		PrintDataSet(&dat)
		t.Fatal()
	}
}
