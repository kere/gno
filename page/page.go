package page

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
	Init(req *http.Request, ps httprouter.Params)

	Auth() (require, isok bool, redirectURL string)
	Prepare() error
	Render(io.Writer) error
}

// Page class
type Page struct {
	Title string
	Name  string
	Dir   string

	Theme string
	Head  []string
	JS    []string
	CSS   []string

	Top    map[string]interface{}
	Bottom map[string]interface{}

	Data interface{}

	Request *http.Request
	Params  httprouter.Params
}

// AddHead page head
func (p *Page) AddHead(items ...string) {
	p.Head = append(p.Head, items...)
}

// AddJS page js
func (p *Page) AddJS(items ...string) {
	p.JS = append(p.JS, items...)
}

// AddCSS page css
func (p *Page) AddCSS(items ...string) {
	p.CSS = append(p.CSS, items...)
}

// AddTop page top template
func (p *Page) AddTop(name string, data interface{}) {
	if p.Top == nil {
		p.Top = make(map[string]interface{}, 0)
	}
	p.Top[name] = data
}

// AddBottom page bottom template
func (p *Page) AddBottom(name string, data interface{}) {
	if p.Bottom == nil {
		p.Bottom = make(map[string]interface{}, 0)
	}
	p.Bottom[name] = data
}

// Init page
func (p *Page) Init(req *http.Request, ps httprouter.Params) {
	p.Request = req
	p.Params = ps
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
	lyt := layout.NewPage(p.Title)

	name := filepath.Join(p.Dir, p.Name)
	for _, v := range p.Head {
		lyt.Head.AddHeadItem(render.NewHead(v))
	}

	for _, v := range p.JS {
		lyt.Head.AddJS(render.NewJS(filepath.Join(name, v)))
	}

	for _, v := range p.CSS {
		lyt.Head.AddCSS(render.NewCSS(filepath.Join(name, v), p.Theme))
	}

	for k, v := range p.Top {
		lyt.AddTop(filepath.Join(p.Dir, k), v)
	}

	lyt.AddBody(name+".htm", p.Data)

	for k, v := range p.Bottom {
		lyt.AddBottom(filepath.Join(p.Dir, k), v)
	}

	return lyt.Render(w)
}
