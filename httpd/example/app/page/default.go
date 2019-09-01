package page

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/httpd/render"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
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
	// requirejs
	data := make(map[string]string, 0)
	data["defer"] = ""
	data["async"] = "true"

	data["data-main"] = httpd.Site.AssetsURL + util.PathToURL("/js/", httpd.RunMode+"/page", d.D.Dir, d.D.Name)
	data["src"] = "/assets/js/require.js"

	// css := render.NewCSS("https://cdn.jsdelivr.net/npm/element-ui@2.11.1/lib/theme-chalk/index.css")
	d.D.CSS = []render.IRenderWith{render.NewCSS("default.css")}
	d.D.JS = []render.IRenderWith{render.NewJS("vue.js"), render.Script("", data)}

	// d.D.Head = []render.IRender{}
	d.D.Top = []render.IRender{render.NewTemplate("_header.htm")}

	// d.D.CacheOption.PageMode = httpd.CacheModePagePath
	// d.D.CacheOption.Store = httpd.CacheStoreNone
	d.D.CacheOption.Store = httpd.CacheStoreMem

	d.D.Bottom = []render.IRender{render.NewTemplate("_bottom.htm")}
	return d
}

// Page page
func (d *Default) Page(ctx *fasthttp.RequestCtx) error {
	// time.Sleep(3 * time.Second)
	return nil
}

//
// // Auth page
// func (d *Default) Auth(ctx *fasthttp.RequestCtx) error {
// 	return nil
// }
