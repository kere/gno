package httpd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/kere/gno/db"
	"github.com/kere/gno/layout"
	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
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

	// Location *time.Location
	Router *fasthttprouter.Router

	AssetsURL string
	ErrorURL  string
	LoginURL  string

	JSVer  string
	CSSVer string

	Secret string
	Nonce  string

	Log  *log.Logger
	PID  string
	Lang []byte

	Pool int
}

// Init Server
func Init() *SiteServer {
	config = conf.Load(ConfigName)

	a := config.GetConf("site")
	s := &SiteServer{
		Listen: a.DefaultString("listen", ":8080"),
		Router: fasthttprouter.New(),
	}

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

	// RunMode
	RunMode = a.DefaultString("mode", "dev")

	// AssetsURL
	s.AssetsURL = a.DefaultString("assets_url", "")

	// ErrorURL
	s.ErrorURL = a.DefaultString("error_url", "/error")
	// LoginURL
	s.LoginURL = a.DefaultString("login_url", "/login")

	// Secret
	s.Secret = a.DefaultString("secret", "")
	s.Nonce = fmt.Sprint(time.Now().Unix())

	// PID
	s.PID = a.DefaultString("pid", "")
	// Pool
	s.Pool = a.DefaultInt("pool", 20)
	// Lang
	s.Lang = []byte(a.DefaultString("lang", "zh"))
	initPool(s.Pool)

	Site = s

	// DB
	if config.IsSet("db") {
		db.Init("app", config.GetConf("db"))
	}

	if config.IsSet("cache") {
		cache.Init(config.GetConf("cache"))
		db.SetCache(cache.CurrentCache())
	}

	err := os.MkdirAll(filepath.Dir(cacheFileStoreDir), os.ModeDir)
	if err != nil {
		panic(err)
	}

	return s
}

// Start server listen
func (s *SiteServer) Start() {
	if layout.RunMode == "dev" {
		s.Router.NotFound = fasthttp.FSHandler(WEBROOT, 0)
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
	svr := &fasthttp.Server{
		ErrorHandler: s.ErrorHandler,
		Handler:      s.Router.Handler,
	}

	util.ListenSignal(func(sign os.Signal) {
		if sign == os.Interrupt {
			if s.PID != "" {
				os.Remove(s.PID)
			}

			if err := svr.Shutdown(); err != nil {
				fmt.Println(err)
			}
		}
		os.Exit(0)
	})

	go func() {
		<-quitChan
	}()

	if err := svr.ListenAndServe(s.Listen); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// ErrorHandler do error
func (s *SiteServer) ErrorHandler(ctx *fasthttp.RequestCtx, err error) {

}
