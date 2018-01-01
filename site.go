package goo

import (
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/goo/layout"
	"github.com/kere/goo/libs/conf"
	"github.com/kere/goo/libs/log"
	"github.com/kere/goo/page"
	"github.com/kere/goo/render"
)

var (
	// ConfigName config file name
	ConfigName = "app/app.conf"

	// SiteServer svr
	SiteServer *Site
)

// Site class
type Site struct {
	Addr       string
	EnableGzip bool
	Location   *time.Location
	Router     *httprouter.Router

	Log *log.Logger
}

// Init goo
func Init() *Site {
	config := conf.Load(ConfigName)

	a := config.GetConf("site")
	s := &Site{
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

	render.JSVersion = a.DefaultString("js_version", "")
	render.CSSVersion = a.DefaultString("css_version", "")
	render.TemplateLeftDelim = a.DefaultString("template_left_delim", "")
	render.TemplateLeftDelim = a.DefaultString("template_right_delim", "")

	layout.RunMode = a.DefaultString("mode", "dev")

	render.RunMode = layout.RunMode
	render.AssetsURL = a.DefaultString("assets_url", "")
	SiteServer = s

	return s
}

// Start server listen
func (s *Site) Start() {
	if layout.RunMode == "dev" {
		s.Router.ServeFiles("/assets/*filepath", http.Dir("webroot/assets"))
	}

	http.ListenAndServe(s.Addr, s.Router)
}

// RegistGet router
func (s *Site) RegistGet(rule string, factory func() page.IPage) {

	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()

		p.Init(req, ps)
		if isReq, isOK, url := p.Auth(); isReq && !isOK {
			if url != "" {
				http.Redirect(rw, req, url, http.StatusNotAcceptable)
			}
			return
		} else if isReq && isOK && url != "" {
			http.Redirect(rw, req, url, http.StatusOK)
		}

		if err := p.Prepare(); err != nil {
			doError(rw, err)
			s.Log.Alert(err)
		}

		if err := p.Render(rw); err != nil {
			doError(rw, err)
			s.Log.Alert(err)
		}

	})
}

func doError(rw http.ResponseWriter, err error) {
	rw.Write([]byte(err.Error()))
}
