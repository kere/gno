package httpd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/kere/gno/db"
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

func init() {
	var err error
	HomeDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
}

// SiteData class
type SiteData struct {
	Lang       string
	AssetsURL  string
	JSVersion  string
	CSSVersion string

	Secret string
	Nonce  string

	ErrorURL string
	LoginURL string
}

// SiteServer class
type SiteServer struct {
	Name   string
	Listen string

	// Location *time.Location
	Router   *fasthttprouter.Router
	SiteData *SiteData

	Log *log.Logger

	PID        string
	Lang       []byte
	Server     *fasthttp.Server
	ConfigName string
	C          conf.Configuration

	AllowFilesHandle bool
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
	TemplateLeftDelim = a.DefaultString("template_left_delim", "")
	TemplateRightDelim = a.DefaultString("template_right_delim", "")

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
		SiteData:   &SiteData{},
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
	site.SiteData.ErrorURL = a.DefaultString("error_url", "/error")
	// LoginURL
	site.SiteData.LoginURL = a.DefaultString("login_url", "/login")

	// Secret
	site.SiteData.Secret = a.DefaultString("secret", "")
	site.SiteData.Nonce = fmt.Sprint(time.Now().Unix())

	// PID
	site.PID = a.DefaultString("pid", "")
	// Pool
	// s.Pool = a.DefaultInt("pool", 0)
	// initPool(s.Pool)

	// Lang
	site.Lang = []byte(a.DefaultString("lang", "zh"))

	site.Name = a.DefaultString("name", "httpd")
	site.SiteData.AssetsURL = a.DefaultString("assets_url", "/assets")

	// JsVersion CSSVersion
	site.SiteData.JSVersion = a.DefaultString("js_version", "")
	site.SiteData.CSSVersion = a.DefaultString("css_version", "")

	site.AllowFilesHandle = a.DefaultBool("allow_files_handle", true)

	site.Server = &fasthttp.Server{
		ErrorHandler: site.ErrorHandler,
		Handler:      site.Router.Handler,
	}

	return site
}

// Start server listen
func (s *SiteServer) Start() {
	if s.AllowFilesHandle {
		s.Router.NotFound = fasthttp.FSHandler(WEBROOT, 0)
	}
	if !DisablePageCache {
		os.MkdirAll(cacheFileStoreDir, os.ModeDir)
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
