package gno

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/layout"
	"github.com/kere/gno/render"
)

var (
// templateExt = ".htm"
)

// IPage interface
type IPage interface {
	GetName() string
	GetDir() string
	GetTheme() string
	GetResponseWriter() http.ResponseWriter
	GetRequest() *http.Request
	GetParams() httprouter.Params

	AddHead(src string)
	AddJS(filename string)
	AddCSS(filename string)
	AddTop(filename string, data interface{})
	AddBottom(filename string, data interface{})
	AddScript(position, src string, data map[string]string)

	Init(method string, w http.ResponseWriter, req *http.Request, ps httprouter.Params)
	Auth() (require, isok bool, redirectURL string, err error)
	// Build() error
	Prepare() error
	Render() error
	SetCookie(name, value string, age int, path, domain string, httpOnly bool)
}

// Page class
type Page struct {
	Title string
	Name  string
	Dir   string

	Theme string
	Lang  string

	Head []render.IRender
	JS   []render.IRender
	CSS  []render.IRender

	Top    []render.IRender
	Bottom []render.IRender

	Data interface{}

	Method         string
	Request        *http.Request
	Params         httprouter.Params
	ResponseWriter http.ResponseWriter

	// Layout     *layout.Page
	JSPosition string
}

// Init page
func (p *Page) Init(method string, w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	p.Method = method
	p.Request = req
	p.Params = params
	p.ResponseWriter = w
	p.Lang = Site.Lang
	p.Theme = Site.Theme
}

// GetName value
func (p *Page) GetName() string {
	return p.Name
}

// GetDir value
func (p *Page) GetDir() string {
	return p.Dir
}

// GetTheme value
func (p *Page) GetTheme() string {
	return p.Theme
}

// GetRequest value
func (p *Page) GetRequest() *http.Request {
	return p.Request
}

// GetResponseWriter value
func (p *Page) GetResponseWriter() http.ResponseWriter {
	return p.ResponseWriter
}

// GetParams value
func (p *Page) GetParams() httprouter.Params {
	return p.Params
}

// AddHead head
func (p *Page) AddHead(src string) {
	r := render.NewHead(src)
	p.Head = append(p.Head, r)
}

// AddCSS css file
func (p *Page) AddCSS(filename string) {
	r := render.NewCSS(filename, p.Theme)
	p.CSS = append(p.CSS, r)
}

// AddJS js file
func (p *Page) AddJS(filename string) {
	r := render.NewJS(filename)
	p.JS = append(p.JS, r)
}

// AddTop add a top render
func (p *Page) AddTop(filename string, data interface{}) {
	r := render.NewTemplate(filename)
	r.Data = data
	p.Top = append(p.Top, r)
}

// AddBottom add a bottom render
func (p *Page) AddBottom(filename string, data interface{}) {
	r := render.NewTemplate(filename)
	r.Data = data
	p.Bottom = append(p.Bottom, r)
}

// AddScript add a bottom render
// position: top, bottom
func (p *Page) AddScript(position, src string, data map[string]string) {
	str := "<script"

	var s string
	if len(data) > 0 {
		s = " "
		for k, v := range data {
			s += k + "=\"" + v + "\" "
		}
	}

	str += s + ">" + src + "</script>"
	switch position {
	case "bottom":
		p.Bottom = append(p.Bottom, render.NewString(str))
	case "top":
		p.Top = append(p.Top, render.NewString(str))
	default:
		p.Head = append(p.Head, render.NewString(str))
	}
}

// Auth page auth
// if require is true then do auth
func (p *Page) Auth() (require, isok bool, redirectURL string, err error) {
	return false, false, "", nil
}

// Prepare page
func (p *Page) Prepare() error {
	return nil
}

// SetCookie cookie
func (p *Page) SetCookie(name, value string, age int, path, domain string, httpOnly bool) {
	SetCookie(p.ResponseWriter, name, value, age, path, domain, httpOnly)
}

// SetCookie f
func SetCookie(w http.ResponseWriter, name, value string, age int, path, domain string, httpOnly bool) {
	var expires time.Time

	if age != 0 {
		expires = time.Unix(time.Now().Unix()+int64(age), 0)
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		Expires:  expires,
		HttpOnly: httpOnly,
	}

	http.SetCookie(w, cookie)
}

// Render page
func (p *Page) Render() error {
	lyt := layout.NewPage()
	// lyt := p.Layout
	lyt.Lang = p.Lang

	lyt.Head.Theme = p.Theme
	lyt.Head.Title = p.Title
	lyt.Head.CSSRenders = p.CSS
	lyt.Head.JSRenders = p.JS
	lyt.Head.HeadItems = p.Head
	lyt.Head.JSPosition = p.JSPosition

	name := filepath.Join(p.Dir, p.Name)

	lyt.Top = p.Top

	if name != "" {
		lyt.AddBody(name+".htm", p.Data)
	}

	lyt.Bottom = p.Bottom

	// p.Layout = nil
	p.Params = nil
	p.Request = nil

	err := lyt.Render(p.ResponseWriter)
	p.ResponseWriter = nil
	return err
}
