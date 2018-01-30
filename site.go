package gno

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/db"
	"github.com/kere/gno/layout"
	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/render"
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
	Listen   string
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

	Log *log.Logger
	PID string
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

	// ----------- log -------------
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
	s.PID = a.DefaultString("pid", "")

	Site = s

	// DB
	if config.IsSet("db") {
		db.Init("app", config.GetConf("db"))
	}

	return s
}

// Start server listen
func (s *SiteServer) Start() {
	if layout.RunMode == "dev" {
		// s.Router.ServeFiles("/assets/*filepath", http.Dir("webroot/assets"))
		s.Router.NotFound = http.FileServer(http.Dir("webroot"))
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

// RegistGet router
func (s *SiteServer) RegistGet(rule string, factory func() IPage) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("GET", rw, req, ps)
		err := doPageHandle(p, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doPageError(s.ErrorURL, err, rw, req)
	})
}

// RegistPost router
func (s *SiteServer) RegistPost(rule string, factory func() IPage) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("POST", rw, req, ps)
		err := doPageHandle(p, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doPageError(s.ErrorURL, err, rw, req)
	})
}

// RegistAPI api router
func (s *SiteServer) RegistAPI(rule string, factory func() IWebAPI) {
	f := func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		webapi := factory()
		webapi.Init(rw, req, ps)
		err := doAPIHandle(webapi, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doAPIError(err, rw)
	}

	s.Router.POST(rule, f)
	s.Router.GET(rule, f)
}

// RegistPostAPI api router
func (s *SiteServer) RegistPostAPI(rule string, factory func() IWebAPI) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		webapi := factory()
		webapi.Init(rw, req, ps)
		err := doAPIHandle(webapi, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doAPIError(err, rw)
	})
}

// RegistGetAPI api router
func (s *SiteServer) RegistGetAPI(rule string, factory func() IWebAPI) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		webapi := factory()
		webapi.Init(rw, req, ps)
		err := doAPIHandle(webapi, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doAPIError(err, rw)
	})
}

func doPageHandle(p IPage, rw http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	isReq, isOK, urlstr, err := p.Auth()
	if isReq && !isOK {
		if urlstr != "" {
			u, _ := url.Parse(urlstr)
			if u.RawQuery == "" {
				u.RawQuery = "msg=" + url.PathEscape(err.Error())
			} else {
				u.RawQuery += "&msg=" + url.PathEscape(err.Error())
			}

			http.Redirect(rw, req, u.String(), http.StatusSeeOther)
		}
		return nil
	}

	err = p.Prepare()
	if err != nil {
		return err
	}

	if p.GetName() == "" {
		return nil
	}

	return p.Render()
}

func doAPIHandle(webapi IWebAPI, rw http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	if isReq, err := webapi.Auth(); isReq && err != nil {
		return err
	}

	data, err := webapi.Exec()
	if err != nil {
		return err
	}

	return webapi.Reply(data)
}

func doPageError(errorURL string, err error, rw http.ResponseWriter, req *http.Request) {
	log.App.Warn(err)
	if errorURL == "" {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	// ErrorURL redirect to
	http.Redirect(rw, req, errorURL+"?msg="+err.Error(), http.StatusSeeOther)
}

func doAPIError(err error, rw http.ResponseWriter) {
	log.App.Warn(err)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte(err.Error()))
}
