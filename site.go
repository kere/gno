package gno

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/db"
	"github.com/kere/gno/layout"
	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/render"
)

const (
	// WEBROOT string
	WEBROOT = "webroot"
	// ModePro pro
	ModePro = "pro"
	// ModeDev dev
	ModeDev = "dev"
)

var (
	// ConfigName config file name
	ConfigName = "app/app.conf"
	// RunMode dev pro
	RunMode = "dev"
	// Site svr
	Site *SiteServer

	config conf.Configuration

	quitChan = make(chan os.Signal)
)

// SiteServer class
type SiteServer struct {
	Listen string

	Location *time.Location
	Router   *httprouter.Router

	ErrorURL   string
	JSVersion  string
	CSSVersion string
	AssetsURL  string

	Lang  string
	Theme string

	Secret string
	Salt   string

	Log     *log.Logger
	PID     string
	HomeDir string
}

// GetConfig return config
func GetConfig() conf.Configuration {
	return config
}

// Init goo
func Init() *SiteServer {
	config = conf.Load(ConfigName)

	a := config.GetConf("site")
	s := &SiteServer{
		Listen: a.DefaultString("listen", ":8080"),
		Router: httprouter.New()}

	//  log
	if config.IsSet("log") {
		l := config.GetConf("log")

		if _, err := os.Stat("var/log"); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir("var/log", os.ModePerm)
			}
		}

		log.Init("var/log/", l.DefaultString("logname", "app"), l.DefaultString("logstore", "stdout"), l.DefaultString("level", "info"))

	} else {
		log.Init("", "app", "stdout", "")
	}
	s.Log = log.Get("app")

	// ------- time zone --------
	if a.IsSet("timezone") {
		zone := a.GetString("timezone")
		loc, err := time.LoadLocation(zone)
		if err != nil {
			panic(err)
		}
		s.Location = loc
	}
	pool = NewPool(a.DefaultInt("pool", 200))

	// JsVersion CSSVersion
	s.JSVersion = a.DefaultString("js_version", "")
	render.JSVersion = s.JSVersion
	s.CSSVersion = a.DefaultString("css_version", "")
	render.CSSVersion = s.CSSVersion

	// Template Delim
	render.TemplateLeftDelim = a.DefaultString("template_left_delim", "")
	render.TemplateRightDelim = a.DefaultString("template_right_delim", "")

	// RunMode
	RunMode = a.DefaultString("mode", "dev")
	layout.RunMode = RunMode

	// AssetsURL
	s.AssetsURL = a.DefaultString("assets_url", "")
	render.AssetsURL = s.AssetsURL

	// ErrorURL
	s.ErrorURL = a.DefaultString("error_url", "")

	// Secret
	s.Secret = a.DefaultString("secret", "")
	s.Salt = fmt.Sprint(time.Now().Unix())

	// Lang
	s.Lang = a.DefaultString("lang", "en")
	// Theme
	s.Theme = a.DefaultString("theme", "")
	// PID
	s.PID = a.DefaultString("pid", "")

	dir := filepath.Dir(ConfigName)
	s.HomeDir, _ = filepath.Abs(filepath.Join(dir, ".."))

	Site = s

	// DB
	if config.IsSet("db") {
		db.Init("app", config.GetConf("db"))
	}

	if config.IsSet("cache") {
		cache.Init(config.GetConf("cache"))
	}

	return s
}

// Start server listen
func (s *SiteServer) Start() {
	if layout.RunMode == "dev" {
		// s.Router.ServeFiles("/assets/*filepath", http.Dir("webroot/assets"))
		s.Router.NotFound = http.FileServer(http.Dir(WEBROOT))
	}

	if s.PID != "" {
		f, err := os.OpenFile(s.PID, os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		f.WriteString(fmt.Sprint(os.Getpid()))
		f.Close()
	}

	fmt.Println("RunMode:", RunMode)
	fmt.Println("Listen:", s.Listen)
	// go http.ListenAndServe(s.Listen, s.Router)
	server := &http.Server{Addr: s.Listen, Handler: s.Router}

	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-quitChan
		if s.PID != "" {
			os.Remove(s.PID)
		}

		if err := server.Close(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// DoError report error
func DoError(err error) error {
	Site.Log.Error(err).Stack()
	return err
}

// DoRecover dillwith panic
func DoRecover() {
	err := recover()
	if err != nil {
		Site.Log.Error(err).Stack()
	}
}
