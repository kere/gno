package log

import (
	"bytes"
	"fmt"
	"io"
	golog "log"
	"os"
	"path/filepath"
	"runtime/debug"
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

	l.Logger.SetPrefix(prefix)
	m = append(m, "\n")
	l.Logger.Println(m...)

	return l
	// if PrintStackLevel < loglevel {
	// 	return l
	// }
	//
	// return l.Stack()
}

func (l *Logger) writelogf(prefix string, loglevel int, format string, m []interface{}) *Logger {
	if l.level < loglevel {
		return l
	}

	l.Logger.SetPrefix(prefix)
	l.Logger.Printf(format+"\n\n", m...)

	return l
	// if PrintStackLevel < loglevel {
	// 	return l
	// }
	//
	// return l.Stack()
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
	// l.Logger.SetPrefix("[emerg]")
	l.Logger.Print("[emerg]")
	l.Logger.Println(m...)
	l.Write("emergency exit(4)")
	os.Exit(4)
}

// Emergf log
func (l *Logger) Emergf(format string, m ...interface{}) {
	// l.Logger.SetPrefix("[emerg]")

	l.Logger.Printf("[emerg]"+format, m...)
	l.Write("\nemergency exit(4)")
	os.Exit(4)
}

// Alert log
func (l *Logger) Alert(m ...interface{}) *Logger {
	return l.writelog("[alert]", LogAlert, m...)
}

// Alertf log
func (l *Logger) Alertf(format string, m ...interface{}) *Logger {
	return l.writelogf("[alert]", LogAlert, format, m)
}

// Crit log
func (l *Logger) Crit(m ...interface{}) *Logger {
	return l.writelog("[crit]", LogCrit, m...)
}

// Critf log
func (l *Logger) Critf(format string, m ...interface{}) *Logger {
	return l.writelog("[crit]", LogCrit, format, m)
}

func (l *Logger) Error(m ...interface{}) *Logger {
	return l.writelog("[err]", LogErr, m...)
}

// Errorf log
func (l *Logger) Errorf(format string, m ...interface{}) *Logger {
	return l.writelogf("[err]", LogErr, format, m)
}

// Warn log
func (l *Logger) Warn(m ...interface{}) *Logger {
	return l.writelog("[warn]", LogWarning, m...)
}

// Warnf log
func (l *Logger) Warnf(format string, m ...interface{}) *Logger {
	return l.writelogf("[warn]", LogWarning, format, m)
}

// Notice log
func (l *Logger) Notice(m ...interface{}) *Logger {
	return l.writelog("[notice]", LogNotice, m...)
}

// Noticef log
func (l *Logger) Noticef(format string, m ...interface{}) *Logger {
	return l.writelogf("[notice]", LogNotice, format, m)
}

// Info log
func (l *Logger) Info(m ...interface{}) *Logger {
	return l.writelog("[info]", LogInfo, m...)
}

// Infof log
func (l *Logger) Infof(format string, m ...interface{}) *Logger {
	return l.writelogf("[info]", LogInfo, format, m)
}

// Debug log
func (l *Logger) Debug(m ...interface{}) *Logger {
	return l.writelog("[debug]", LogDebug, m...)
}

// Debugf log
func (l *Logger) Debugf(format string, m ...interface{}) *Logger {
	return l.writelogf("[debug]", LogDebug, format, m)
}

// Sql log
func (l *Logger) Sql(sqlstr []byte, args []interface{}) *Logger {
	if l.level < LogSQL {
		return l
	}

	for i, item := range args {
		switch item.(type) {
		case []byte, []int8:
			args[i] = string(item.([]byte))
		}
	}

	return l.writelog("[sql]", LogSQL, string(sqlstr), args)
}
