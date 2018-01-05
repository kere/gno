package gno

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
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
	// AssetsURL js css img
	AssetsURL = ""
	// Site svr
	Site *SiteServer
	// JSVersion js ?v=001
	JSVersion = ""
	// CSSVersion css ?v=001
	CSSVersion = ""

	config conf.Configuration
)

// SiteServer class
type SiteServer struct {
	Addr       string
	EnableGzip bool
	Location   *time.Location
	Router     *httprouter.Router

	Log *log.Logger
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
		Addr:       a.DefaultString("addr", ":8080"),
		EnableGzip: a.DefaultBool("gzip", false),
		Router:     httprouter.New()}

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
	JSVersion = a.DefaultString("js_version", "")
	render.JSVersion = JSVersion
	CSSVersion = a.DefaultString("css_version", "")
	render.CSSVersion = CSSVersion

	// Template Delim
	render.TemplateLeftDelim = a.DefaultString("template_left_delim", "")
	render.TemplateRightDelim = a.DefaultString("template_right_delim", "")

	// RunMode
	RunMode = a.DefaultString("mode", "dev")
	layout.RunMode = RunMode

	// AssetsURL
	AssetsURL = a.DefaultString("assets_url", "")
	render.AssetsURL = AssetsURL

	Site = s

	return s
}

// Start server listen
func (s *SiteServer) Start() {
	if layout.RunMode == "dev" {
		// s.Router.ServeFiles("/assets/*filepath", http.Dir("webroot/assets"))
		s.Router.NotFound = http.FileServer(http.Dir("webroot"))
	}

	fmt.Println("RunMode:", RunMode)
	fmt.Println("Listen:", s.Addr)
	http.ListenAndServe(s.Addr, s.Router)
}

// RegistGet router
func (s *SiteServer) RegistGet(rule string, factory func() IPage) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("GET", rw, req, ps)
		doPageHandle(p, rw, req, ps)
	})
}

// RegistPost router
func (s *SiteServer) RegistPost(rule string, factory func() IPage) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("POST", rw, req, ps)
		doPageHandle(p, rw, req, ps)
	})
}

// RegistAPI api router
func (s *SiteServer) RegistAPI(rule string, factory func() IWebAPI) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		webapi := factory()
		webapi.Init(rw, req, ps)
		doAPIHandle(webapi, rw, req, ps)
	})
}

func doPageHandle(p IPage, rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if isReq, isOK, url := p.Auth(); isReq && !isOK {
		if url != "" {
			http.Redirect(rw, req, url, http.StatusSeeOther)
		}

		return

	} else if isReq && isOK && url != "" {
		http.Redirect(rw, req, url, http.StatusSeeOther)
	}

	err := p.Prepare()
	if err != nil {
		doError(rw, err)
		return
	}

	if p.GetName() == "" {
		return
	}

	err = p.Render()
	if err != nil {
		doError(rw, err)
		return
	}
}

func doAPIHandle(webapi IWebAPI, rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if isReq, isOK := webapi.Auth(); isReq && !isOK {
		doError(rw, errors.New("auth failed"))
		return
	}

	data, err := webapi.Exec()
	if err != nil {
		doError(rw, err)
		return
	}

	webapi.Reply(data)
}

func doError(rw http.ResponseWriter, err error) {
	log.App.Warn(err)
	rw.Write([]byte(err.Error()))
}
