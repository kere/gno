package gno

import (
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/layout"
	"github.com/kere/gno/render"
)

// IPage interface
type IPage interface {
	GetName() string
	GetDir() string
	GetTheme() string
	GetResponseWriter() http.ResponseWriter
	GetRequest() *http.Request
	GetParams() httprouter.Params

	GetCacheMode() int
	GetExpires() int
	SetPageCache(mode, expires int)

	AddHead(src string)
	AddJS(filename string)
	AddHeadJS(filename string)
	SetJSPosition(pos string)
	AddCSS(filename string)
	AddTop(filename string, data interface{})
	AddBottom(filename string, data interface{})
	AddScript(position, src string, data map[string]string)

	AddTopRender(r render.IRender)
	AddBottomRender(r render.IRender)

	Init(method string, w http.ResponseWriter, req *http.Request, ps httprouter.Params)
	Auth() (require, isok bool, redirectURL string, err error)
	// Build() error
	Prepare() error
	Render(io.Writer) error
	SetCookie(name, value string, age int, path, domain string, httpOnly bool)

	// AddAfter(PageExec)  //add page after exec
	// AddBefore(PageExec) // add page before exec
	// RunBefore()
	RunAfter()
}

// PageExec run page exec
type PageExec func(IPage)

// Page class
type Page struct {
	Title string
	Name  string
	Dir   string

	CacheMode int // 0: no cache   1: cache page   2:cache by url path
	Expires   int //page expires

	HTML string

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

	AfterExecs  []PageExec
	BeforeExecs []PageExec
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

// SetJSPosition set
func (p *Page) SetJSPosition(pos string) {
	p.JSPosition = pos
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

// AddHeadJS js file
func (p *Page) AddHeadJS(filename string) {
	r := render.NewJS(filename)
	p.Head = append(p.Head, r)
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

// AddTopRender add body render template
func (p *Page) AddTopRender(r render.IRender) {
	p.Top = append(p.Top, r)
}

// AddBottomRender add body render template
func (p *Page) AddBottomRender(r render.IRender) {
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
	} else {
		s = " type=\"text/javascript\""
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
	SetCookie(p.GetResponseWriter(), name, value, age, path, domain, httpOnly)
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
func (p *Page) Render(w io.Writer) error {
	lyt := layout.NewPage()
	lyt.Lang = p.Lang

	lyt.Head.Theme = p.Theme
	lyt.Head.Title = p.Title
	lyt.Head.CSSRenders = p.CSS
	lyt.Head.JSRenders = p.JS
	lyt.Head.HeadItems = p.Head
	lyt.Head.JSPosition = p.JSPosition

	lyt.Top = p.Top

	if p.HTML != "" {
		lyt.AddBodyRender(render.NewString(p.HTML))
	} else if p.Name != "" {
		name := filepath.Join(p.Dir, p.Name)
		lyt.AddBody(name+".htm", p.Data)
	}

	lyt.Bottom = p.Bottom

	// p.Layout = nil
	err := lyt.Render(w)
	return err
}

// // AddBefore page
// func (p *Page) AddBefore(e PageExec) {
// 	p.BeforeExecs = append(p.BeforeExecs, e)
// }

// // AddAfter page
// func (p *Page) AddAfter(e PageExec) {
// 	p.AfterExecs = append(p.AfterExecs, e)
// }

// RunAfter page
func (p *Page) RunAfter() {
	// l := len(p.AfterExecs)
	// for i := 0; i < l; i++ {
	// 	p.AfterExecs[i](p)
	// }

	p.Params = nil
	p.Request = nil
	p.ResponseWriter = nil
}

// // RunBefore page
// func (p *Page) RunBefore() {
// 	l := len(p.BeforeExecs)
// 	for i := 0; i < l; i++ {
// 		p.BeforeExecs[i](p)
// 	}
// }
