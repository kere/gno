# Gno lib
## i18n package useage
```go
i18n.SetSourceRoot("./source")
i18n.Load("en-US", "page/message")
//path: ./source/en-US/page/message.json
T := i18n.TransFunc("en-US", "page/message")
m := T("title1", "hello!!!")
fmt.Println(m)
```

## log package useage
```go
log.Init("folder", []string{"web", "sql"}, "dev")
log.Level = 10
log.Use("sql")
log.App.Error("error", "code: ", 2)
// log error and print debug info
log.App.Error("error", "something ", "here").Stack()
// auto print debug info
log.PrintStackLevel = 5
```
## conf package useage
config file
```conf
[app]
key=value
key1=value1
key2=value2
key3=value3
key4=value4
```
```go
configuration := conf.Load("app/app.conf")
conf := conf["app"]
boolval  := conf.IsSet("key")
string   := conf.GetString("key1")
intval   := conf.GetInt("key2")
floatval := conf.GetFloat("key3")
boolval  := conf.GetBool("key4")
```
