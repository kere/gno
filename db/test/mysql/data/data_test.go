package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kere/gno/db"
	"github.com/kere/gno/db/drivers"
	"github.com/kere/gno/libs/conf"
	_ "github.com/lib/pq"
)

func Test_data_initdb(t *testing.T) {
	cf := conf.Conf{}
	cf["driver"] = "mysql"
	cf["dbname"] = "testdb"
	cf["host"] = "127.0.0.1"
	cf["port"] = "3306"
	cf["user"] = "root"
	cf["password"] = "123123"
	cf["level"] = "info"

	db.New("app", cf)
	if db.Current() == nil {
		t.Error("current database is nil")
	}

}

func Test_drivers_postgres(t *testing.T) {
	pgsql := &drivers.Postgres{}
	src := []byte("{{0,3}}")
	v := make([][2]int64, 0)
	err := pgsql.ParseNumberSlice(src, &v)
	if err != nil {
		t.Fatal(err)
	}

	if v[0][0] != int64(0) || v[0][1] != int64(3) {
		t.Fatalf("postgres parse slice failed")
	}
}

func Test_postgres_row(t *testing.T) {
	row := db.DataRow{}
	row["id"] = 1
	row["name"] = "tomas"
	row["number"] = 22.22
	row["data"] = `{"text":"this is a message."}`
	row["int_arr"] = "[2,3,4]"
	row["num_arr"] = "[2.1,3.2,4.3]"
	row["string_arr"] = `a,b,c`
	row["created_at"] = time.Now()

	if a := row.GetInt64Slice("int_arr"); a[0] != 2 {
		t.Fatalf("row GetInt64Slice failed, get vaule is %d", a[0])
	}
	/*if a := row.GetInt64Slice("num_arr"); a[0] != 2 {
		t.Fatalf("row GetInt64Slice failed, get vaule is %s", a)
	}*/
	if a := row.Floats("num_arr"); a[0] != 2.1 {
		t.Fatalf("row Floats failed, get vaule is %.1f", a[0])
	}
	if a := row.GetStringSlice("string_arr"); a[0] != "a" {
		t.Fatalf("row GetStringSlice failed, get vaule is %s. the first value is %s", a, a[0])
	}
}

type UserVO1 struct {
	db.BaseVO
	Id      int64     `json:"id" skip:"insert"`
	Name    string    `json:"name"`
	Age     int       `json:"age"`
	Data    string    `json:"data"`
	Created time.Time `json:"created_at" autotime:"true" skip:"update"`
}

func Test_rowdata2struct(t *testing.T) {
	row := db.DataRow{}
	row["id"] = 1
	row["name"] = "tomas"
	row["age"] = 22
	row["data"] = "this is a message."
	row["created_at"] = time.Now()

	vo := &UserVO1{}
	err := row.CopyToStruct(vo)
	if err != nil {
		t.Fatal(err)
	}
	if vo.Id != row.GetInt64("id") {
		t.Fatalf("vo.Id %d is not equal row[id] %d", vo.Id, row.GetInt("id"))
	}

	if vo.Name != row.GetString("name") {
		t.Fatalf("vo.Name %s is not equal row[name] %s", vo.Name, row.GetString("name"))
	}

	if vo.Age != row.GetInt("age") {
		t.Fatalf("vo.Age %d is not equal row[age] %d", vo.Age, row.GetInt("age"))
	}

	if vo.Created.Unix() != row.GetTime("created_at").Unix() {
		t.Fatalf("vo.Created %s is not equal row[created_at] %s", vo.Created, row.GetString("created_at"))
	}
}

type UserVO2 struct {
	db.BaseVO
	Data interface{} `json:"data"`
}

func Test_rowdata2struct_typeto_interface_type(t *testing.T) {
	row := db.DataRow{}
	row["data"] = `{"text":"this is a message."}`

	vo := &UserVO2{}
	err := row.CopyToStruct(vo)
	if err != nil {
		t.Error(err.Error())
	}

	if reflect.TypeOf(vo.Data).String() != "map[string]interface {}" {
		t.Fatalf("type error: vo.Data [%s] is not map[string]interface {}", reflect.TypeOf(vo.Data).String())
	} else {
		data := db.DataRow(vo.Data.(map[string]interface{}))
		if data.GetString("text") != "this is a message." {
			t.Fatalf("data[text] `%s` is not equal `this is a message.`", data.GetString("text"))
		}
	}

	m := make(map[string]interface{}, 0)
	m["text"] = "this is a message."

	row = db.DataRow{}
	row["data"] = m

	vo = &UserVO2{}
	err = row.CopyToStruct(vo)
	if err != nil {
		t.Error(err)
	}

	if vo.Data == nil {
		t.Fatalf("vo.Data is nil")
	}

	data := db.DataRow(vo.Data.(map[string]interface{}))
	if data.GetString("text") != "this is a message." {
		t.Fatalf("data[text] `%s` is not equal `this is a message.`", data.GetString("text"))
	}

}

type UserVO3 struct {
	db.BaseVO
	Data db.DataRow `json:"data"`
}

func Test_rowdata2struct_typeto_datarow(t *testing.T) {
	row := db.DataRow{}
	row["data"] = `{"text":"this is a message."}`

	vo := &UserVO3{}
	err := row.CopyToStruct(vo)
	if err != nil {
		t.Error(err.Error())
	}

	if fmt.Sprint(vo.Data["text"]) != "this is a message." {
		t.Fatalf("data[text] `%s` is not equal `this is a message.`", fmt.Sprint(vo.Data["text"]))
	}

	// start with map data
	row = db.DataRow{}
	m := make(map[string]interface{}, 0)
	m["text"] = "this is a message."
	row["data"] = m

	vo = &UserVO3{}
	err = row.CopyToStruct(vo)
	if err != nil {
		t.Error(err.Error())
	}

	if fmt.Sprint(vo.Data["text"]) != "this is a message." {
		t.Fatalf("data[text] `%s` is not equal `this is a message.`", fmt.Sprint(vo.Data["text"]))
	}
}

type UserVO4 struct {
	db.BaseVO
	Data  map[string]interface{} `json:"data"`
	Users [][2]int64             `json:"users"`
}

func Test_rowdata2struct_typeto_map(t *testing.T) {
	row := db.DataRow{}
	row["data"] = `{"text":"this is a message."}`
	row["users"] = "[[0,3]]"

	vo := &UserVO4{}
	err := row.CopyToStruct(vo)
	if err != nil {
		t.Error(err.Error())
	}

	if fmt.Sprint(vo.Data["text"]) != "this is a message." {
		t.Fatalf("data.Text `%s` is not equal `this is a message.`", fmt.Sprint(vo.Data["text"]))
	}

	if vo.Users[0][0] != int64(0) || vo.Users[0][1] != int64(3) {
		t.Fatalf("[][2]int64 value is wrong, it convert failed")
	}
}

type UserVO5 struct {
	db.BaseVO
	Data *SubData `json:"data"`
}

type SubData struct {
	Text string `json:"text"`
}

func Test_rowdata2struct_typeto_struct(t *testing.T) {
	row := db.DataRow{}
	row["data"] = `{"text":"this is a message."}`

	vo := &UserVO5{}
	err := row.CopyToStruct(vo)
	if err != nil {
		t.Error(err.Error())
	}

	if vo.Data.Text != "this is a message." {
		t.Fatalf("sub struct data.Text `%s` is not equal `this is a message.`", vo.Data.Text)
	}
}

func Test_struct2datarow(t *testing.T) {
	vo := GetUserVO_a1()
	sv := db.NewStructConvert(vo)

	row := sv.Struct2DataRow("update")

	v, isOk := row["name"]
	if !isOk || v != "kere" {
		t.Fatalf("name is %s, expect kere", v)
	}

	if _, isOk = row["age2"]; isOk {
		t.Fatalf("update skip age2")
	}

	if _, isOk = row["Name3"]; isOk {
		t.Fatalf("skipempty Name3")
	}
}

type UserVO_a1 struct {
	db.BaseVO
	Id      int64                  `json:"id" skip:"all"`
	Name    string                 `json:"name"`
	Age     int                    `json:"age"`
	Name2   string                 `json:"name2" skip:"insert"`
	Age2    int                    `json:"age2" skip:"update"`
	Name3   string                 `json:"name3" skipempty:"all"`
	Data    interface{}            `json:"data"`
	Users1  [][2]int64             `json:"users1"`
	Users2  []int64                `json:"users2"`
	Names   []string               `json:"names"`
	SubData *SubData               `json:"subdata"`
	DataMap map[string]interface{} `json:"datamap"`
	DataRow db.DataRow             `json:"datarow"`
	Created time.Time              `json:"created_at" autotime:"true" skip:"update"`
}

func GetUserVO_a1() *UserVO_a1 {
	vo := &UserVO_a1{}
	vo.Id = int64(1)
	vo.Name = "kere"
	vo.Age = 22

	vo.Name2 = "kere2"
	vo.Age2 = 23

	vo.Users1 = make([][2]int64, 1)
	vo.Users1[0] = [2]int64{0, 3}

	vo.Users2 = []int64{8, 9}
	vo.Names = []string{"tom1", "tom2"}

	subdata := &SubData{"this is a message."}
	vo.SubData = subdata

	datarow := db.DataRow{}
	datarow["text"] = "this is a message."
	vo.DataRow = datarow

	m := make(map[string]interface{}, 0)
	m["text"] = "this is a message."
	vo.DataMap = m

	vo.Data = m

	vo.Created = time.Now()

	return vo
}

type AskForLeaveFromData struct {
	db.BaseVO
	EmptyStr string      `json:"empty_str"`
	Content  string      `json:"content"`
	BeginAt  time.Time   `json:"begin_at"`
	EndAt    time.Time   `json:"end_at"`
	Users1   [][2]int64  `json:"users1"`
	Users2   []int64     `json:"users2"`
	Users3   [2]int64    `json:"users3"`
	StrArr1  []string    `json:"strarr1"`
	StrArr2  [][2]string `json:"strarr2"`
	StrArr3  [2]string   `json:"strarr3"`
}

func (n *AskForLeaveFromData) Table() string {
	return ""
}

func Test_row2struct(t *testing.T) {
	s := `{"content":"this is a message.","empty_str":null,"begin_at":"2014-11-25 16:00:00","end_at":"2014-11-27 16:00:00","users1":[[0,6],[1,5],[2,4]],"users2":[1,3,4],"users3":[1,2],"strarr1":["a","b"],"strarr2":[["a","b"],["a2","b2"]],"strarr3":["a3","b3"]}`
	row := db.DataRow{}
	err := json.Unmarshal([]byte(s), &row)
	if err != nil {
		t.Fatalf("json parse error: %s \nsrc=%s", err.Error(), s)
	}

	vo := &AskForLeaveFromData{}
	row.CopyToStruct(vo)

	if vo.Content != "this is a message." {
		t.Fatalf("v.Content is %s", vo.Content)
	}

	if vo.BeginAt.Day() != 25 {
		t.Fatalf("v.BeginAt is %s", vo.BeginAt)
	}

	if vo.EmptyStr != "" {
		t.Fatalf("empty str is %s", vo.EmptyStr)
	}

	if vo.Users1[0][1] != int64(6) {
		t.Fatalf("vo.Users1[0][1] = %d, is must = 6", vo.Users1[0][1])
	}
	if vo.Users2[1] != int64(3) {
		t.Fatalf("vo.Users2[1] = %d, is must = 3", vo.Users2[1])
	}
	if vo.Users3[0] != int64(1) {
		t.Fatalf("vo.Users3[0] = %d, is must = 1", vo.Users3[0])
	}

	if vo.StrArr1[1] != "b" {
		t.Fatalf("vo.StrArr1[1] = %s, it must = b", vo.StrArr1[1])
	}
	if vo.StrArr2[1][1] != "b2" {
		t.Fatalf("vo.StrArr2[1][1] = %s, it must = b2", vo.StrArr2[1][1])
	}
	if vo.StrArr3[1] != "b3" {
		t.Fatalf("vo.StrArr3[1] = %s, it must = b3", vo.StrArr2[1][1])
	}
}

func Test_struct_parser(t *testing.T) {
	vo := GetUserVO_a1()

	// insert mode
	convert := db.NewStructConvert(vo)

	// update mode
	keys, values, _ := convert.KeyValueList("update")
	// id is skiped when insert mode.
	for i, k := range keys {
		if string(k) == `"name"=?` && values[i] != "kere" {
			t.Fatalf("name is %s, it must equal kere", values[i])
		} else if string(k) == `"age"=?` && fmt.Sprint(values[i]) != "22" {
			t.Fatalf("age is %s, it must equal 22", fmt.Sprint(values[i]))
		}
	}

	// id is skiped when insert mode.
	keys, values, _ = convert.KeyValueList("insert")

	for i, k := range keys {
		switch string(k) {
		case `"id"`:
			t.Fatalf("id 'skip' tab invalid")

		case `"name"`:
			if fmt.Sprint(values[i]) != vo.Name {
				t.Fatalf("name value is %s, it must equal %s", fmt.Sprint(values[0]), vo.Name)
			}

		case `"data"`, `"subdata"`, `"datarow"`, `"datamap"`:
			if fmt.Sprintf("%s", values[i]) != `{"text":"this is a message."}` {
				t.Fatalf("%s is not a json string. %s", k, values[i])
			}

		case `"users1"`:
			if fmt.Sprintf("%s", values[i]) != "{{0,3}}" {
				t.Fatalf("users1 value is %s , it must equal '{{0,3}}", values[i])
			}

		case `"users2"`:
			if fmt.Sprintf("%s", values[i]) != "{8,9}" {
				t.Fatalf("users2 value is %s , it must equal '{8,9}", values[i])
			}

		case `"names"`:
			if fmt.Sprintf("%s", values[i]) != "{'tom1','tom2'}" {
				t.Fatalf("names value is %s , it must equal {'tom1','tom2'}", values[i])
			}

		}
	}

}

//--------------------------- Benchmark ---------------------------

func Benchmark_copy_to_struct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		row := db.DataRow{}
		row["id"] = 1
		row["name"] = "tomas"
		row["age"] = 22
		row["data"] = "this is a message."
		row["created_at"] = time.Now()

		vo := &UserVO1{}
		err := row.CopyToStruct(vo)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_struct_parser(b *testing.B) {
	vo := GetUserVO_a1()
	convert := db.NewStructConvert(vo)
	for i := 0; i < b.N; i++ {
		convert.KeyValueList("insert")
	}
}

func Benchmark_struct_parser2(b *testing.B) {
	vo := GetUserVO_a1()
	convert := db.NewStructConvert(vo)
	for i := 0; i < b.N; i++ {
		convert.KeyValueList("update")
	}
}

//---------------------------
