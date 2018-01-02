package goo

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/goo/layout"
	"github.com/kere/goo/render"
)

var (
// templateExt = ".htm"
)

// IPage interface
type IPage interface {
	AddHead(items ...string)
	AddJS(items ...string)
	AddCSS(items ...string)
	Init(method string, req *http.Request, ps httprouter.Params)
	Auth() (require, isok bool, redirectURL string)
	Build() error
	Prepare() error
	Render(io.Writer) error
}

// Page class
type Page struct {
	Title string
	Name  string
	Dir   string

	Theme string
	Head  []render.IRender
	JS    []render.IRender
	CSS   []render.IRender

	Top    []render.IRender
	Bottom []render.IRender

	Data interface{}

	Method  string
	Request *http.Request
	Params  httprouter.Params
	Layout  *layout.Page
}

// Init page
func (p *Page) Init(method string, req *http.Request, ps httprouter.Params) {
	p.Method = method
	p.Request = req
	p.Params = ps

	p.Layout = layout.NewPage()
}

// Build page
func (p *Page) Build() error {
	return nil
}

// AddHead head
func (p *Page) AddHead(src string) {
	r := render.NewHead(src)
	p.Head = append(p.Head, r)
}

// AddCSS css file
func (p *Page) AddCSS(filename, theme string) {
	r := render.NewCSS(filepath.Join(p.Dir, p.Name, filename), p.Theme)
	p.CSS = append(p.CSS, r)
}

// AddJS js file
func (p *Page) AddJS(filename string) {
	r := render.NewJS(filepath.Join(p.Dir, p.Name, filename))
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

// AddBottomScript add a bottom render
func (p *Page) AddBottomScript(src string, data map[string]string) {
	str := "<script "
	for k, v := range data {
		str += k + "=\"" + v + "\" "
	}
	str += ">" + src + "</script>"
	p.Bottom = append(p.Bottom, render.NewString(str))
}

// Auth page auth
// if require is true then do auth
func (p *Page) Auth() (require, isok bool, redirectURL string) {
	return false, false, ""
}

// Prepare page css
func (p *Page) Prepare() error {
	return nil
}

// Render page css
func (p *Page) Render(w io.Writer) error {
	lyt := p.Layout
	p.Layout = nil

	lyt.Head.Title = p.Title

	name := filepath.Join(p.Dir, p.Name)

	lyt.AddBody(name+".htm", p.Data)

	return lyt.Render(w)
}
