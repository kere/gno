package test

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kere/gno/db"
	"github.com/kere/gno/libs/conf"
)

var testTableName = "users"

func Test_init(t *testing.T) {
	cf := conf.Conf{}
	cf["driver"] = "mysql"
	cf["dbname"] = "testdb"
	cf["addr"] = "127.0.0.1:3306"
	cf["user"] = "root"
	cf["password"] = "123123"
	cf["level"] = "all"

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

	// str := []byte(fmt.Sprintf("TRUNCATE TABLE %s;", testTableName))
	// _, err := db.Current().Exec(db.NewSqlState(str, nil))
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
	_, err := db.Current().ExecStr(fmt.Sprintf("drop table `testdb`.`%s`;", testTableName))
	if err != nil {
		t.Log(err)
	}

	sqlstr := "	CREATE TABLE IF NOT EXISTS `testdb`.`" + testTableName + "` ("
	sqlstr += "  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,"
	sqlstr += "  `name` VARCHAR(45) NULL,"
	sqlstr += "  `age` INT NULL,"
	sqlstr += "  `context` TEXT NULL,"
	sqlstr += "  `data_json` JSON NULL,"
	sqlstr += "  `names` JSON NULL,"
	sqlstr += "  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,"
	sqlstr += "  PRIMARY KEY (`id`)) "
	sqlstr += "ENGINE = InnoDB;"

	_, err = db.Current().ExecStr(sqlstr)
	if err != nil {
		t.Fatal(err)
		return
	}
}

type DataJson struct {
	Date      time.Time
	Code      string
	Price     float64
	IntSlice  []int
	IntString []string
}

type RowVO struct {
	db.BaseVO
	Id       int64     `json:"id" skip:"all"`
	Name     string    `json:"name" skipempty:"update"`
	Age      int       `json:"age"`
	Context  string    `json:"context"`
	Names    []string  `json:"names"`
	DataJson *DataJson `json:"data_json"`
	Created  time.Time `json:"created_at" skip:"insert" autotime:"true"`
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
	names := []string{"a", "b", "c"}
	datajson := &DataJson{time.Now(), "00001", 21.12, []int{1, 2, 3}, names}
	row1 := db.DataRow{"name": "kere1", "age": 21, "names": names, "data_json": datajson}
	row2 := db.DataRow{"name": "kere2", "age": 22, "names": names, "data_json": datajson}
	row3 := db.DataRow{"name": "kere3", "age": 23, "names": names, "data_json": datajson}
	row4 := db.DataRow{"name": "kere4", "age": 24, "names": names, "data_json": datajson}
	row5 := db.DataRow{"name": "kere5", "age": 25, "names": names, "data_json": datajson}

	rows := db.DataSet{row2, row3, row4}

	_, err := ins.Insert(row1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ins.InsertM(rows)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.NewInsertBuilder(testTableName).IsPrepare(false).Insert(row5)
	if err != nil {
		t.Fatal(err)
	}

	row, err := db.NewQueryBuilder(testTableName).Where("name=?", "kere1").QueryOne()
	if err != nil {
		t.Fatal(err)
	}
	if !row.IsNull("context") {
		t.Fatal("context isnull failed")
	}

	_, err = db.NewQueryBuilder(testTableName).IsPrepare(false).Where("name=?", "kere1").QueryOne()
	if err != nil {
		t.Fatal(err)
	}

	// struct --------------------------
	// ins.Begin()
	// userVO := &RowVO{Name: "Struct1", Age: 22, Context: "Struct Context1", Names: names, DataJson: datajson}
	// _, err = ins.Insert(userVO)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// lastID := ins.LastInsertId("id")
	// if lastID != 5 {
	// 	t.Errorf("insert id is %d", lastID)
	// }
	// ins.End()

	userVO := &RowVO{Name: "tom", Age: 33, Context: "Struct Context2", Names: names, DataJson: datajson}
	_, err = ins.Insert(userVO)
	if err != nil {
		t.Fatal(err)
		return
	}

	row, err = db.NewQueryBuilder(testTableName).Where("name=?", "tom").QueryOne()
	if err != nil {
		t.Fatal(err)
		return
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

	// test skipempty
	vo.Name = ""
	err = vo.Update(vo.Id)
	if err != nil {
		t.Fatal(err)
	}

	vo, err = findOne("tom")
	if err != nil {
		t.Fatal(err)
	}

	if vo.Name == "" {
		t.Fatal("skipempty failed")
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

	dataset, err := db.ScanRows(rows)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		db.NewStructConvert(cls).DataSet2Struct(dataset)
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
	if d == nil {
		return nil, fmt.Errorf("%s not found.", n)
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
