package page

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/libs/util"
	"github.com/kere/gno/render"
	"github.com/valyala/fasthttp"
)

// Default page class
type Default struct {
	name, dir string
	title     string

	head   []render.IRender
	css    []render.IRender
	js     []render.IRender
	top    []render.IRender
	bottom []render.IRender
}

// NewDefault func
func NewDefault() *Default {
	d := &Default{}
	d.title = "Default Page"
	d.name = "default"
	d.dir = ""

	d.css = []render.IRender{render.NewCSS("default.css", "")}

	d.head = []render.IRender{}
	d.top = []render.IRender{render.NewTemplate("_header.htm")}

	// requirejs
	data := make(map[string]string, 0)
	data["defer"] = ""
	data["async"] = "true"

	data["data-main"] = httpd.Site.AssetsURL + util.PathToURL("/assets/js/", httpd.RunMode+"/page", d.Dir(), d.Name())
	data["src"] = "/assets/js/require.js"
	d.bottom = []render.IRender{render.NewTemplate("_bottom.htm"), render.Script("", data)}

	return d
}

// Title page
func (d *Default) Title() string {
	return d.title
}

// Lang page
func (d *Default) Lang() string {
	return ""
}

// Dir page
func (d *Default) Dir() string {
	return d.dir
}

// Name page
func (d *Default) Name() string {
	return d.name
}

// CSS page
func (d *Default) CSS() []render.IRender {
	return d.css
}

// JS page
func (d *Default) JS() []render.IRender {
	return d.js
}

// Head page
func (d *Default) Head() []render.IRender {
	return d.head
}

// Top page
func (d *Default) Top() []render.IRender {
	return d.top
}

// Body page
func (d *Default) Body() []render.IRender {
	return nil
}

// Bottom page
func (d *Default) Bottom() []render.IRender {
	return d.bottom
}

// CacheOption page
func (d *Default) CacheOption() (mode, headExpires, expires int) {
	return httpd.CacheModeNone, 30, 30
}

// Page page
func (d *Default) Page(ctx *fasthttp.RequestCtx) error {
	return nil
}

// Auth page
func (d *Default) Auth(ctx *fasthttp.RequestCtx) (string, error) {
	return "", nil
}
