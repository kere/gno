package page

import (
	"github.com/kere/gno/httpd"
	"github.com/valyala/fasthttp"
)

// Default page class
type Default struct {
	httpd.P
}

// NewDefault func
func NewDefault() *Default {
	d := &Default{}
	d.Init("Default Page", "Default", "")

	// requirejs

	d.PA.CSS = []httpd.IRenderA{httpd.NewCSS("default.css")}

	d.PA.JS = []httpd.IRenderA{httpd.NewJS("vue.js"), httpd.RequireJS(&d.PA, "/assets/js/require.js", nil)}

	// d.PA.Head = []httpd.IRender{}
	d.PA.Top = []httpd.IRender{httpd.NewTemplate("_header.htm")}

	// d.PA.CacheOption.PageMode = httpd.CacheModePagePath
	// d.PA.CacheOption.Store = httpd.CacheStoreNone
	d.PA.CacheOption.Store = httpd.CacheStoreMem

	d.PA.Bottom = []httpd.IRender{httpd.NewTemplate("_bottom.htm")}
	return d
}

// Page page
func (d *Default) Page(ctx *fasthttp.RequestCtx) (*httpd.PageData, error) {
	// time.Sleep(3 * time.Second)
	return nil, nil
}

//
// // Auth page
// func (d *Default) Auth(ctx *fasthttp.RequestCtx) error {
// 	return nil
// }
