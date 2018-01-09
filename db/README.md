# db package useage

**begin use db package, you need install and import db driver by your self, like this:**
```go
import(_ "github.com/mattn/go-sqlite3")
```

```go
type DataSet []interface{}
type DataRow map[string]interface{}
```
#### 1. query
```go
q := db.NewQueryBuilder("users")   //q := (&db.QueryBuilder{}).Table("users")
q.Table("users").Where("name=? and age=?", "tom", 22)
ds, err := q.Query() //return DataSet
r, err :=q.QueryOne() //return DataRow
```

#### 2. query result with struct   
```go
type UserVO struct{
	ID int64 `json:"id" skip:"all"`
	Name string `json:"name", skip:"update"`
	Age int64 `json:"age" skipempty:"all"`
	Updated time.Time `json:"updated_at" autotime:"true" skip:"insert"`
	Created time.Time `json:"created_at" autotime:"true" skip:"update"`
}
// select name,age,created_at from Usres
q := db.NewQueryBuilder("users") //q := (&db.QueryBuilder{}).Table("users")
```
```go
q.Struct(&UserVO{}) // or
q.Struct((*UserVO)(nil))
r, err := q.FindOne() //return interface{}.(*UserVO)
items, err := q.Find() //return []interface{}
```
```go
q.Struct(UserVO{})
r, err := q.FindOne() //return interface{}.(UserVO)
```
#### 3. UpdateBuilder
```go
u := db.NewUpdateBuilder("users").Where("id=?", 1)
rowData := db.RowData{}
rowData["name"] = "toms"
u.Update(rowData)
```
or
```go
row := db.DataRow{}
row["name"] = "tom"
row["age"] = 23
u.Update(row)
```
or update with data struct
```go
rowVO := &UserVO{Name:"toms", Age:23}
// rowVO := &UserVO{}
// row.CopyToStruct(rowVO)
u.Update(rowVO)
```

#### 4. DeleteBuilder
```go
d:=db.NewDeleteBuilder("users").Where("id=?", 1).Delete()
```

#### 5. CountBuilder
```go
count := db.NewCountBuilder("users").Count()
```
#### 6. ExistsBuilder
```go
e, err := db.NewExistsBuilder("users").Where("id=?", 1).Exists()
```

#### 7. db.Query() and db.QueryX()
```go
db.Query("select * from Users") //return []DataRow
```
```go
db.Find(&UserVO{}, "select * from Users") //return []interface{}

```
#### 8. db.Tx{}
```go
tx := &db.Tx{}
tx.Begin()
insert := db.NewInsertBuilder("users")
row := UserVO{"tom", 22}
insert.TxInsert(tx, row) // or tx.Exec(insert.SqlState(row))
lastId, err := tx.LastInsertId("users", "id")
delete := db.NewDeleteBuilder("users").Where("nick=?", "lucy")
delete.TxDelete(tx) // or tx.Exec(delete.SqlState())
tx.End()
```

#### 9. db cache
```go
q.Cache().Query() // cache result, it will user db.DefaultQueryCacheExpire.
// or
q.CacheExpire(300).Query() // cache result
```
clear cache
```go
q.ClearCache()
```
**sqlite3**
```go
conf:=make(map[string]string)
conf["driver"] = "sqlite3"
conf["file"] = "./app.db"
// or
conf["driver"] = "sqlite3"
conf["connect"] = "file=./app.db"
```
**postgres**
```go
conf:=make(map[string]string)
conf["driver"] = "postgres"
conf["dbname"] = "mydb"
conf["host"] = "127.0.0.1"
conf["port"] = "5432"
conf["user"] = "postgres"
conf["password"] = "123"
// or
conf["connect"] = "dbname=mydb user=postgres password=123 host=127.0.0.1 port=5432 sslmode=disable"
```
**init db pool**
```go
db.Init("app", conf)
db.New("app2", conf2)
db.Use("app2")
```
