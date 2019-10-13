package log

import (
	"bytes"
	"fmt"
	"io"
	golog "log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

const (
	sBBreak = "\n\n"
	backetL = "["
	backetR = "]"
)

type emptyWriter struct{}

func (e *emptyWriter) Write(p []byte) (int, error) {
	return 0, nil
}

// NewEmpty func
func NewEmpty() *Logger {
	l := &Logger{}
	l.Logger = golog.New(&emptyWriter{}, "", 0)
	return l
}

// Logger class
type Logger struct {
	level     int
	LevelName string
	*golog.Logger
}

// NewLogger func
func NewLogger(file, levelStr string) *Logger {
	return (&Logger{}).Init(file, levelStr)
}

// Init func
func (l *Logger) Init(file, levelStr string) *Logger {
	var out io.Writer
	l.level = IntLevel(levelStr)
	if l.level == -100 {
		levelStr = "not found"
	}
	l.LevelName = levelStr

	if file == "" {
		fmt.Println("Logger is write in Stdout")
		fmt.Println("Logger level name:", l.LevelName, " level value:", l.level)
		out = os.Stdout
	} else {
		folder, _ := filepath.Split(file)
		if _, err := os.Stat(folder); err != nil {
			if os.IsNotExist(err) {
				os.MkdirAll(folder, os.ModePerm)
			}
		}

		fmt.Println("Logger is write in *.log file: " + file)
		fmt.Println("Logger level name:", l.LevelName, " level value:", l.level)
		var err error
		out, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			golog.Fatalln(err)
		}
		// os.Chmod(file, 0664)
	}
	l.Logger = golog.New(out, "", 0)
	l.Logger.SetFlags(golog.Ldate | golog.Ltime) // | golog.Lshortfile

	return l
}

// SetLevel func
func (l *Logger) SetLevel(s string) *Logger {
	if l == nil {
		return l
	}
	l.LevelName = s
	l.level = IntLevel(s)
	return l
}

// SetPrefix func
func (l *Logger) SetPrefix(prefix string) *Logger {
	l.Logger.SetPrefix(prefix)
	return l
}

func (l *Logger) Write(m ...interface{}) *Logger {
	l.Logger.Println(m...)
	return l
}

func (l *Logger) writelog(prefix string, loglevel int, m ...interface{}) *Logger {
	if l.level < loglevel {
		return l
	}

	var s strings.Builder
	s.WriteString(backetL)
	s.WriteString(prefix)
	s.WriteString(backetR)
	n := len(m)
	arr := make([]interface{}, 2+n)
	arr[0] = s.String()
	arr[n+1] = "\n"
	for i := 0; i < n; i++ {
		arr[i+1] = m[i]
	}

	l.Logger.Println(arr...)
	return l
}

func (l *Logger) writelogf(prefix string, loglevel int, format string, m []interface{}) *Logger {
	if l.level < loglevel {
		return l
	}

	// l.Logger.SetPrefix(prefix)
	l.Logger.Printf(format+sBBreak, m...)

	return l
}

var bBreak = []byte{'\n'}
var bSearchStrGno = []byte("github.com/kere/gno")

// Stack print
func (l *Logger) Stack() *Logger {
	count := 0
	arr := bytes.Split(debug.Stack(), bBreak)
	arr = arr[3:]
	for i := range arr {
		if bytes.Index(arr[i], bSearchStrGno) > -1 {
			count++
		} else {
			break
		}
	}

	if count > 0 {
		arr = arr[count:]
	}

	l.Write("\n", string(bytes.Join(arr, bBreak)))
	return l
}

// Emerg log
func (l *Logger) Emerg(m ...interface{}) {
	l.Logger.Print(LogEmergStr)
	l.Logger.Println(m...)
	l.Write("emergency exit(4)")
	l.Stack()
	os.Exit(4)
}

// Emergf log
func (l *Logger) Emergf(format string, m ...interface{}) {
	l.Logger.Printf(LogEmergStr+" "+format, m...)
	l.Write("\nemergency exit(4)")
	l.Stack()
	os.Exit(4)
}

// Alert log
func (l *Logger) Alert(m ...interface{}) *Logger {
	l.Stack()
	return l.writelog(LogAlertStr, LogAlert, m...)
}

// Alertf log
func (l *Logger) Alertf(format string, m ...interface{}) *Logger {
	l.Stack()
	return l.writelogf(LogAlertStr, LogAlert, format, m)
}

// Crit log
func (l *Logger) Crit(m ...interface{}) *Logger {
	return l.writelog(LogCritStr, LogCrit, m...)
}

// Critf log
func (l *Logger) Critf(format string, m ...interface{}) *Logger {
	return l.writelog(LogCritStr, LogCrit, format, m)
}

func (l *Logger) Error(m ...interface{}) *Logger {
	return l.writelog(LogErrStr, LogErr, m...)
}

// Errorf log
func (l *Logger) Errorf(format string, m ...interface{}) *Logger {
	return l.writelogf(LogErrStr, LogErr, format, m)
}

// Warn log
func (l *Logger) Warn(m ...interface{}) *Logger {
	return l.writelog(LogWarnStr, LogWarn, m...)
}

// Warnf log
func (l *Logger) Warnf(format string, m ...interface{}) *Logger {
	return l.writelogf(LogWarnStr, LogWarn, format, m)
}

// Notice log
func (l *Logger) Notice(m ...interface{}) *Logger {
	return l.writelog(LogNoticeStr, LogNotice, m...)
}

// Noticef log
func (l *Logger) Noticef(format string, m ...interface{}) *Logger {
	return l.writelogf(LogNoticeStr, LogNotice, format, m)
}

// Info log
func (l *Logger) Info(m ...interface{}) *Logger {
	return l.writelog(LogInfoStr, LogInfo, m...)
}

// Infof log
func (l *Logger) Infof(format string, m ...interface{}) *Logger {
	return l.writelogf(LogInfoStr, LogInfo, format, m)
}

// Debug log
func (l *Logger) Debug(m ...interface{}) *Logger {
	return l.writelog(LogDebugStr, LogDebug, m...)
}

// Debugf log
func (l *Logger) Debugf(format string, m ...interface{}) *Logger {
	return l.writelogf(LogDebugStr, LogDebug, format, m)
}

// Sql log
func (l *Logger) Sql(dbname, sqlstr string, args []interface{}) *Logger {
	return l.writelog(LogSQLStr, LogSQL, backetL+dbname+backetR, sqlstr, args)
}
