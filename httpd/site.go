package httpd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/kere/gno/db"
	"github.com/kere/gno/httpd/render"
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

	// HomeDir of config
	HomeDir = ""
	// RunMode dev pro
	RunMode = "dev"
	// Site svr
	Site *SiteServer

	config conf.Configuration

	quitChan = make(chan os.Signal)
)

// SiteServer class
type SiteServer struct {
	Name   string
	Listen string

	// Location *time.Location
	Router *fasthttprouter.Router

	// AssetsURL string
	ErrorURL string
	LoginURL string

	Secret string
	Nonce  string

	Log *log.Logger

	PID  string
	Lang []byte

	// Timeout       time.Duration
	MaxConnsPerIP int
	Concurrency   int //连接并发数
	// Pool int
}

// GetConfig return Configuration
func GetConfig() *conf.Configuration {
	return &config
}

// Init Server
func Init(name string) {
	ConfigName = name
	config = conf.Load(name)

	dir := filepath.Dir(name)
	HomeDir, _ = filepath.Abs(filepath.Join(dir, ".."))

	a := config.GetConf("site")
	Site = &SiteServer{
		Listen: a.DefaultString("listen", ":8080"),
		Router: fasthttprouter.New(),
	}

	//  log
	if config.IsSet("log") {
		l := config.GetConf("log")

		// if _, err := os.Stat("var/log"); err != nil {
		// 	if os.IsNotExist(err) {
		// 		os.Mkdir("var/log", os.ModePerm)
		// 	}
		// }
		log.Init("var/log/", l.DefaultString("logname", "app"), l.DefaultString("logstore", log.StoreTypeStd), l.DefaultString("level", "info"))

	} else {
		log.Init("", "app", log.StoreTypeStd, "")
	}
	Site.Log = log.Get("app")

	// RunMode
	RunMode = a.DefaultString("mode", "dev")

	// ErrorURL
	Site.ErrorURL = a.DefaultString("error_url", "/error")
	// LoginURL
	Site.LoginURL = a.DefaultString("login_url", "/login")

	// Secret
	Site.Secret = a.DefaultString("secret", "")
	Site.Nonce = fmt.Sprint(time.Now().Unix())

	// PID
	Site.PID = a.DefaultString("pid", "")
	// Pool
	// s.Pool = a.DefaultInt("pool", 0)
	// initPool(s.Pool)

	// Lang
	Site.Lang = []byte(a.DefaultString("lang", "zh"))

	// Timeout
	// s.Timeout = time.Second * time.Duration(a.DefaultInt("timeout", 2))

	// MaxConnsPerIP
	Site.MaxConnsPerIP = a.DefaultInt("max_conns_per_ip", 0)
	// Concurrency
	Site.Concurrency = a.DefaultInt("concurrency", 2048)
	Site.Name = a.DefaultString("name", "httpd")

	// DB
	if config.IsSet("db") {
		db.Init("app", config.GetConf("db"))
	}

	if config.IsSet("cache") {
		cache.Init(config.GetConf("cache"))
		db.SetCache(cache.CurrentCache())
	}

	// AssetsURL
	render.AssetsURL = a.DefaultString("assets_url", "/assets")
	// JsVersion CSSVersion
	render.JSVersion = a.DefaultBytes("js_version", nil)
	render.CSSVersion = a.DefaultBytes("css_version", nil)
	// Template Delim
	render.TemplateLeftDelim = a.DefaultString("template_left_delim", "")
	render.TemplateRightDelim = a.DefaultString("template_right_delim", "")
}

// Start server listen
func (s *SiteServer) Start() {
	if RunMode == "dev" {
		s.Router.NotFound = fasthttp.FSHandler(WEBROOT, 0)
	}

	if _, err := os.Stat(cacheFileStoreDir); os.IsNotExist(err) {
		os.Mkdir(cacheFileStoreDir, os.ModePerm)
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
		MaxConnsPerIP: s.MaxConnsPerIP,
		Concurrency:   s.Concurrency,
		ErrorHandler:  s.ErrorHandler,
		Handler:       s.Router.Handler,
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
