package page

import (
	"github.com/kere/gno/httpd"
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
	d.D.Init("Default Page", "Default", "")

	// requirejs
	data := make(map[string]string, 0)
	data["defer"] = ""
	data["async"] = "true"

	data["data-main"] = httpd.Site.SiteData.AssetsURL + util.PathToURL("/js/", httpd.RunMode+"/page", d.D.Dir, d.D.Name)
	data["src"] = "/assets/js/require.js"

	// css := httpd.NewCSS("https://cdn.jsdelivr.net/npm/element-ui@2.11.1/lib/theme-chalk/index.css")
	d.D.CSS = []httpd.IRenderWith{httpd.NewCSS("default.css")}

	d.D.JS = []httpd.IRenderWith{httpd.NewJSr("vue.js"), httpd.Script("", data)}

	// d.D.Head = []httpd.IRender{}
	d.D.Top = []httpd.IRender{httpd.NewTemplate("_header.htm")}

	// d.D.CacheOption.PageMode = httpd.CacheModePagePath
	// d.D.CacheOption.Store = httpd.CacheStoreNone
	d.D.CacheOption.Store = httpd.CacheStoreMem

	d.D.Bottom = []httpd.IRender{httpd.NewTemplate("_bottom.htm")}
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
