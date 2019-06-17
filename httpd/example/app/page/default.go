package page

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/libs/util"
	"github.com/kere/gno/render"
)

// Default page class
type Default struct {
	httpd.P
}

// NewDefault func
func NewDefault() *Default {
	d := &Default{}
	d.D.Title = []byte("Default Page")
	d.D.Name = "default"
	d.D.Dir = ""

	d.D.CSS = []render.IRender{render.NewCSS("default.css", "")}

	// d.D.Head = []render.IRender{}
	d.D.Top = []render.IRender{render.NewTemplate("_header.htm")}

	// requirejs
	data := make(map[string]string, 0)
	data["defer"] = ""
	data["async"] = "true"

	data["data-main"] = httpd.Site.AssetsURL + util.PathToURL("/assets/js/", httpd.RunMode+"/page", d.D.Dir, d.D.Name)
	data["src"] = "/assets/js/require.js"
	d.D.Bottom = []render.IRender{render.NewTemplate("_bottom.htm"), render.Script("", data)}
	return d
}

// // Page page
// func (d *Default) Page(ctx *fasthttp.RequestCtx) error {
// 	return nil
// }
//
// // Auth page
// func (d *Default) Auth(ctx *fasthttp.RequestCtx) error {
// 	return nil
// }
