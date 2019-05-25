package db

import (
	"os"
	"path/filepath"

	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
)

// NewLogger func
func NewLogger(dbConf map[string]string) *log.Logger {
	conf := conf.Conf(dbConf)
	levelStr := conf.Get("level")
	level := log.IntLevel(levelStr)

	name := conf.Get("logname")
	if name == "" {
		name = "db"
	}

	if level < 0 {
		return log.New("", name, "std", log.LogAlertStr)
	}

	folder := filepath.Join(filepath.Dir(os.Args[0]), "/var/log/")
	return log.New(folder, name, conf.Get("logstore"), levelStr)
}
