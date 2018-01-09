package log

import (
	"fmt"
	golog "log"
	"os"
	"path/filepath"
	"strings"
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
	LogEmerg   = 1  //* system is unusable */
	LogAlert   = 2  //* action must be taken immediately */
	LogCrit    = 3  //* critical conditions */
	LogErr     = 4  //* error conditions */
	LogWarning = 5  //* warning conditions */
	LogNotice  = 6  //* normal but significant condition */
	LogInfo    = 7  //* informational */
	LogDebug   = 8  //* debug-level messages */
	LogSQL     = 8  //* debug-level messages */
	LogAll     = 10 //* debug-level messages */
	ConstNone  = "none"
)

// IntLevel func
func IntLevel(s string) int {
	switch s {
	case "emerg":
		return LogEmerg
	case "alert":
		return LogAlert
	case "crit":
		return LogCrit
	case "err", "error":
		return LogErr
	case "warn", "warning":
		return LogWarning
	case "notice":
		return LogNotice
	case "info":
		return LogInfo
	case "debug":
		return LogDebug
	case "all":
		return LogAll
	case "close":
		return 0
	case ConstNone, "nil", "empty":
		return -1
	}

	return -100
}

var (
	// PrintStackLevel int
	// PrintStackLevel = 4

	// App default log
	App      *Logger
	Location *time.Location
	pool     loggers
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
func Init(folder, names, storeType, levelStr string) {
	level := IntLevel(levelStr)
	if storeType == ConstNone || level < 0 {
		// App = New("", "", ConstNone, ConstNone)
		return
	}

	if names == "" {
		names = "app"
	}

	arr := strings.Split(names, ",")
	for i, name := range arr {
		if i == 0 {
			App = New(folder, name, storeType, levelStr)
		} else {
			New(folder, name, storeType, levelStr)
		}
	}

}

func Get(name string) *Logger {
	return pool[name]
}

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
		name = "app"
	}

	fmt.Println("Prepare logger:", name)

	level := IntLevel(levelStr)
	if level < 0 || storeType == ConstNone {
		fmt.Println("xxx Logger Closed xxx")
		l = NewEmpty()
	} else if storeType == "file" {
		file := filepath.Join(folder, fmt.Sprint(time.Now().Format("20060102"), "-", name, ".log"))
		l = NewLogger(file, levelStr)
	} else {
		l = NewLogger("", levelStr)
	}

	pool[name] = l
	return l
}

func Use(n string) {
	App = pool[n]
}
