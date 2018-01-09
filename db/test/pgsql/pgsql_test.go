package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kere/gno/db"
	"github.com/kere/gno/libs/conf"
	_ "github.com/lib/pq"
)

var testTableName = "pg_test_table"

func Test_initdb(t *testing.T) {
	cf := conf.Conf{}
	cf["driver"] = "postgres"
	cf["dbname"] = "stockdb_test"
	cf["host"] = "127.0.0.1"
	cf["port"] = "5432"
	cf["user"] = "postgres"
	cf["password"] = "123123"
	cf["level"] = "info"

	d := db.New("app", cf)
	if d == nil {
		t.Fatal("create database failed")
		return
	}
	if d.Connection == nil {
		t.Fatal("conn is nil")
		return
	}
	if db.Current() == nil {
		t.Fatal("current database is nil")
		return
	}

	db.Exec(fmt.Sprintf("drop table %s;", testTableName))

	sqlstr := "create table " + testTableName + " (id BIGSERIAL not null,name VARCHAR(200) null, age INT null, context TEXT null, created_at TIMESTAMP WITH TIME ZONE null, string_arr text[] null, constraint PK_PGTESTTABLE primary key (id));"
	_, err := db.Exec(sqlstr)
	if err != nil {
		t.Fatal(err)
		return
	}

}

type RowVO struct {
	db.BaseVO
	Id      int64     `json:"id" skip:"all"`
	Name    string    `json:"name"`
	Age     int       `json:"age"`
	Context string    `json:"context"`
	Created time.Time `json:"created_at" autotime:"true" skip:"update"`
}

func (u *RowVO) Table() string {
	return testTableName
}

type RowModel struct {
	db.BaseModel
}

func NewRowModel() *RowModel {
	m := &RowModel{}
	m.Init(&RowVO{})
	return m
}

func Test_insertm_by_struct(t *testing.T) {
	ins := db.NewInsertBuilder(testTableName)
	rows := db.DataSet{}
	rows = append(rows, db.DataRow{"name": "kere1", "age": 21}, db.DataRow{"name": "kere2", "age": 23}, db.DataRow{"name": "kere3", "age": 24}, db.DataRow{"name": "kere4", "age": 25})

	_, err := ins.InsertM(rows)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_insert_by_struct(t *testing.T) {
	ins := db.NewInsertBuilder(testTableName)
	userVO := &RowVO{Name: "kere", Age: 22}
	_, err := ins.Insert(userVO)
	if err != nil {
		t.Fatal(err)
	}

	userVO = &RowVO{Name: "tom", Age: 22}
	_, err = ins.Insert(userVO)
	if err != nil {
		t.Fatal(err)
	}

	row, err := db.NewQueryBuilder(testTableName).Where("name=?", "tom").QueryOne()
	if err != nil {
		t.Fatal(err)
	}
	if row.GetTime("created_at").Year() < 2010 {
		t.Fatal("autotime fatal.", row.GetTime("created_at"))
	}
}

func Test_update_by_row_vo(t *testing.T) {
	vo, err := findOne("tom")
	if err != nil {
		t.Fatal(err)
	}

	if vo == nil {
		t.Fatalf("the vo is nil")
	}

	vo.Context = "this is a message."
	err = vo.Update(vo.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_find_all(t *testing.T) {
	ds, err := findAll()
	if err != nil {
		t.Fatal(err)
	}

	ecode := ds.Encode()
	if len(ecode) < 1 {
		t.Fatalf("record count is %d", len(ecode))
	}

	if fmt.Sprint(ecode[0][0]) != "id" {
		t.Fatalf("the first row is %s", ecode[0])
	}
}

func Test_query_by_struct(t *testing.T) {
	vo, err := findOne("tom")
	if err != nil {
		t.Fatal(err)
	}
	if vo.Context != "this is a message." {
		t.Fatalf("user context is not martched. this context is %s", vo.Context)
	}
}

func Test_query_by_map(t *testing.T) {
	row, err := queryOne("tom")
	if err != nil {
		t.Fatal(err)
	}
	if row.GetString("context") != "this is a message." {
		t.Fatalf("row context is not martched. this context is %s", row.GetString("context"))
	}
}

type Record []map[string]interface{}

func Test_query(t *testing.T) {
	query := db.NewQueryBuilder(testTableName).Order("id desc")
	d, err := query.Query()
	if err != nil {
		t.Fatal(err)
	}

	if d[0].GetString("name") != "tom" {
		t.Fatal("user nick is not mathched!")
	}
}

func Test_model(t *testing.T) {
	model := NewRowModel()
	row, err := model.QueryByID(int64(1))
	if err != nil {
		t.Fatal(err)
	}

	if row.Empty() {
		t.Fatalf("query one by id 1, get an empty row data")
	}
}

//--------------------------- Benchmark ---------------------------

func Benchmark_query_by_struct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		findOne("tom")
	}
}

func Benchmark_query_by_map(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queryOne("tom")
	}
}

func Benchmark_scanrow_to_struct(b *testing.B) {
	var cls db.IVO
	cls = &RowVO{}
	rows, _ := db.Current().Connection.Connect().Query("select id,nick,created_at,updated_at from users where name='tomas'")

	for i := 0; i < b.N; i++ {
		db.ScanRowsX(cls, rows)
	}
}

func Benchmark_scanrow_to_map(b *testing.B) {
	rows, _ := db.Current().Connection.Connect().Query("select id,nick,created_at,updated_at from users where name='tomas'")

	for i := 0; i < b.N; i++ {
		db.ScanRows(rows)
	}
}

//---------------------------
func findOne(n string) (*RowVO, error) {
	query := db.NewQueryBuilder(testTableName).Where("name=?", n).Struct(&RowVO{})
	d, err := query.FindOne()
	if err != nil {
		return nil, err
	}

	return d.(*RowVO), nil
}

func findAll() (db.VODataSet, error) {
	query := db.NewQueryBuilder(testTableName).Struct(&RowVO{})
	return query.Find()
}

func queryOne(n string) (db.DataRow, error) {
	query := db.NewQueryBuilder(testTableName).Where("name=?", n)
	d, err := query.QueryOne()
	if err != nil {
		return nil, err
	}
	return d, nil
}
