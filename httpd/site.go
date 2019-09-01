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
	// HomeDir of config
	HomeDir = ""
	// RunMode dev pro
	RunMode = ModeDev
	// Site svr
	Site     *SiteServer
	quitChan = make(chan os.Signal)
)

// SiteServer class
type SiteServer struct {
	Name   string
	Listen string

	// Location *time.Location
	Router *fasthttprouter.Router

	ErrorURL string
	LoginURL string

	AssetsURL  string
	JSVersion  string
	CSSVersion string

	Secret string
	Nonce  string

	Log *log.Logger

	PID        string
	Lang       []byte
	Server     *fasthttp.Server
	ConfigName string
	C          conf.Configuration
}

// Init Server
func Init(name string) {
	Site = New(name)

	dir := filepath.Dir(name)
	HomeDir, _ = filepath.Abs(filepath.Join(dir, ".."))

	a := Site.C.GetConf("site")
	// RunMode
	RunMode = a.DefaultString("mode", "dev")

	// Template Delim
	render.TemplateLeftDelim = a.DefaultString("template_left_delim", "")
	render.TemplateRightDelim = a.DefaultString("template_right_delim", "")

	// DB
	if Site.C.IsSet("db") {
		db.Init("app", Site.C.GetConf("db"))
	}

	if Site.C.IsSet("cache") {
		cache.Init(Site.C.GetConf("cache"))
		db.SetCache(cache.CurrentCache())
	}
}

// New Server
func New(name string) *SiteServer {
	c := conf.Load(name)

	a := c.GetConf("site")
	site := &SiteServer{
		ConfigName: name,
		C:          c,
		Listen:     a.DefaultString("listen", ":8080"),
		Router:     fasthttprouter.New(),
	}

	//  log
	if c.IsSet("log") {
		l := c.GetConf("log")
		log.Init("var/log/", l.DefaultString("logname", "app"), l.DefaultString("logstore", log.StoreTypeStd), l.DefaultString("level", "info"))
	} else {
		log.Init("", "app", log.StoreTypeStd, "")
	}

	site.Log = log.Get("app")

	// ErrorURL
	site.ErrorURL = a.DefaultString("error_url", "/error")
	// LoginURL
	site.LoginURL = a.DefaultString("login_url", "/login")

	// Secret
	site.Secret = a.DefaultString("secret", "")
	site.Nonce = fmt.Sprint(time.Now().Unix())

	// PID
	site.PID = a.DefaultString("pid", "")
	// Pool
	// s.Pool = a.DefaultInt("pool", 0)
	// initPool(s.Pool)

	// Lang
	site.Lang = []byte(a.DefaultString("lang", "zh"))

	site.Name = a.DefaultString("name", "httpd")
	site.AssetsURL = a.DefaultString("assets_url", "/assets")

	// JsVersion CSSVersion
	site.JSVersion = a.DefaultString("js_version", "")
	site.CSSVersion = a.DefaultString("css_version", "")

	site.Server = &fasthttp.Server{
		ErrorHandler: site.ErrorHandler,
		Handler:      site.Router.Handler,
	}

	return site
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

	util.ListenSignal(func(sign os.Signal) {
		if sign == os.Interrupt {
			if s.PID != "" {
				os.Remove(s.PID)
			}

			if err := s.Server.Shutdown(); err != nil {
				fmt.Println(err)
			}
		}
		os.Exit(0)
	})

	go func() {
		<-quitChan
	}()

	if err := s.Server.ListenAndServe(s.Listen); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// ErrorHandler do error
func (s *SiteServer) ErrorHandler(ctx *fasthttp.RequestCtx, err error) {

}

var (
	cssURL string
	jsURL  string
)

// CSSRender new
func (s *SiteServer) CSSRender(name string) render.CSS {
	if cssURL == "" {
		cssURL = s.AssetsURL + "/css/"
	}
	return render.NewCSS(cssURL + name)
}

// JSRender new
func (s *SiteServer) JSRender(name string) render.JS {
	if jsURL == "" {
		jsURL = s.AssetsURL + "/js/"
	}
	return render.NewJS(jsURL + name)
}
