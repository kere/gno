package test

import (
	"os"
	"testing"
	"time"

	"github.com/kere/gno/db"
	_ "github.com/mattn/go-sqlite3"
)

func Test_initdb(t *testing.T) {
	os.Remove("app.db")
	d := db.New("app", nil)
	if d == nil {
		t.Error("create database failed")
		return
	}
	if d.Connect() == nil {
		t.Error("conn is nil")
		return
	}
	if db.Current() == nil {
		t.Error("current database is nil")
		return
	}

	sqlstr := "create table users (id integer not null primary key, nick text, created_at TIMESTAMP not null default CURRENT_TIME, updated_at TIMESTAMP not null default CURRENT_TIME)"
	db.Exec(sqlstr)
}

type UserVO struct {
	ID      int64     `json:"id" insert:"false" update:"false"`
	Nick    string    `json:"nick"`
	Created time.Time `json:"created_at" autotime:"true" skip:"update"`
	Updated time.Time `json:"updated_at" autotime:"true" skip:"insert"`
}

func Test_insert_by_struct(t *testing.T) {
	ins := db.NewInsertBuilder("users")
	userVO := &UserVO{Nick: "kere"}
	_, err := ins.Insert(userVO)
	if err != nil {
		toError(err)
	}

	userVO = &UserVO{Nick: "tom"}
	_, err = ins.Insert(userVO)
	if err != nil {
		toError(err)
	}
}

func Test_update_by_struct(t *testing.T) {
	upd := db.NewUpdateBuilder("users").Where("nick=?", "tom")
	userVO := &UserVO{Nick: "tomas"}
	_, err := upd.Update(userVO)
	if err != nil {
		toError(err)
	}
}

func Test_query_by_struct(t *testing.T) {
	data, err := queryByStruct()
	if err != nil {
		toError(err)
	}
	if data.Nick != "tomas" {
		toError("user nick is not martched!!!")
	}
}

func Test_query_by_map(t *testing.T) {
	data, err := queryByMap()
	if err != nil {
		toError(err)
	}
	if data["nick"] != "tomas" {
		toError("user nick is not martched!!!")
	}
}

type Record []map[string]interface{}

func Test_query(t *testing.T) {
	query := db.NewQueryBuilder("users")
	d, err := query.Query()
	if err != nil {
		toError(err)
	}

	if d[0].GetString("nick") != "kere" {
		toError("user nick is not mathched!")
	}
}

//--------------------------- Benchmark ---------------------------

func Benchmark_query_by_struct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queryByStruct()
	}
}

func Benchmark_query_by_map(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queryByMap()
	}
}

func Benchmark_scanrow_to_struct(b *testing.B) {
	var cls interface{}
	cls = &UserVO{}
	rows, _ := db.Current().Conn().Query("select id,nick,created_at,updated_at from users where nick='tomas'")

	for i := 0; i < b.N; i++ {
		db.ScanRowsX(cls, rows)
	}
}

func Benchmark_scanrow_to_map(b *testing.B) {
	rows, _ := db.Current().Conn().Query("select id,nick,created_at,updated_at from users where nick='tomas'")

	for i := 0; i < b.N; i++ {
		db.ScanRows(rows)
	}
}

//---------------------------

func toError(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func queryByStruct() (*UserVO, error) {
	query := db.NewQueryBuilder("users").Where("nick=?", "tomas").Struct(&UserVO{})
	d, err := query.FindOne()
	if err != nil {
		return nil, err
	}
	return d.(*UserVO), nil
}
func queryByMap() (db.DataRow, error) {
	query := db.NewQueryBuilder("users").Where("nick=?", "tomas")
	d, err := query.QueryOne()
	if err != nil {
		return nil, err
	}
	return d, nil
}
