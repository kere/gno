package log

import (
	"fmt"
	golog "log"
	"os"
	"path/filepath"
	"time"
)

// none,nil,empty=-1
// Close=0
// emerg=1
// alert=2
// crit=3
// err,error=4
// warn,warning=5
// notice=6
// info=7
// debug=8
// all=10
// -100=NotFound
const (
	defaultLogName = "app"
	LogEmerg       = 1  //* system is unusable */
	LogAlert       = 2  //* action must be taken immediately */
	LogCrit        = 3  //* critical conditions */
	LogErr         = 4  //* error conditions */
	LogWarn        = 5  //* warning conditions */
	LogNotice      = 6  //* normal but significant condition */
	LogInfo        = 7  //* informational */
	LogDebug       = 8  //* debug-level messages */
	LogSQL         = 8  //* SQL messages */
	LogFmt         = 9  //* Fmt messages */
	LogAll         = 10 //* all messages */

	LogEmergStr  = "emerg"  //* system is unusable */
	LogAlertStr  = "alert"  //* action must be taken immediately */
	LogCritStr   = "crit"   //* critical conditions */
	LogErrStr    = "err"    //* error conditions */
	LogWarnStr   = "warn"   //* warning conditions */
	LogNoticeStr = "notice" //* normal but significant condition */
	LogInfoStr   = "info"   //* informational */
	LogDebugStr  = "debug"  //* debug-level messages */
	LogSQLStr    = "sql"    //* SQL messages */
	LogFmtStr    = "fmt"    //* Fmt messages */
	LogAllStr    = "all"    //* all messages */
	ConstNone    = "none"
)

// IntLevel func
func IntLevel(s string) int {
	switch s {
	case LogEmergStr:
		return LogEmerg
	case LogAlertStr:
		return LogAlert
	case LogCritStr:
		return LogCrit
	case LogErrStr:
		return LogErr
	case LogWarnStr:
		return LogWarn
	case LogNoticeStr:
		return LogNotice
	case LogInfoStr:
		return LogInfo
	case LogDebugStr:
		return LogDebug
	case LogFmtStr:
		return LogFmt
	case LogAllStr:
		return LogAll
	case ConstNone:
		return -1
	}

	return -100
}

var (
	// PrintStackLevel int
	// PrintStackLevel = 4

	// App default log
	App *Logger
	// Location *time.Location
	pool loggers
)

type loggers map[string]*Logger

func init() {
	pool = make(map[string]*Logger)

	App = &Logger{}
	App.SetLevel("all")
	App.Logger = golog.New(os.Stdout, "", 0)
	App.Logger.SetFlags(golog.Ldate | golog.Ltime)
}

// Init func
func Init(folder, name, storeType, levelStr string) {
	level := IntLevel(levelStr)
	if storeType == ConstNone || level < 0 {
		// App = New("", "", ConstNone, ConstNone)
		return
	}

	// if names == "" {
	// 	names = "app"
	// }

	// arr := strings.Split(names, ",")
	// for _, name := range arr {
	// if name == defaultLogName {
	App = New(folder, name, storeType, levelStr)
	// 	} else {
	// 		New(folder, name, storeType, levelStr)
	// 	}
	// }

}

//Get Logger
func Get(name string) *Logger {
	return pool[name]
}

// New logger
func New(folder, name, storeType, levelStr string) *Logger {
	// if _, isOK := pool[name]; isOK {
	// 	return pool[name]
	// }
	//
	var l *Logger
	if storeType == "" {
		storeType = "stdout"
	}

	if name == "" {
		name = defaultLogName
	}

	fmt.Println("Prepare logger:", name)

	level := IntLevel(levelStr)
	if level < 0 || storeType == ConstNone {
		fmt.Println("xxx Logger Closed xxx")
		l = NewEmpty()
	} else if storeType == "filedate" {
		file := filepath.Join(folder, fmt.Sprint(time.Now().Format("20060102"), "-", name, ".log"))
		l = NewLogger(file, levelStr)
	} else if storeType == "file" {
		file := filepath.Join(folder, name+".log")
		l = NewLogger(file, levelStr)
	} else {
		l = NewLogger("", levelStr)
	}

	pool[name] = l
	return l
}

// Use n log
func Use(n string) {
	App = pool[n]
}
