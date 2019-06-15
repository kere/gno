package page

import (
	"fmt"

	"github.com/kere/gno/httpd"
	"github.com/kere/gno/libs/util"
	"github.com/kere/gno/render"
	"github.com/valyala/fasthttp"
)

// Default page class
type Default struct {
	httpd.PageBase
}

// NewDefault func
func NewDefault() *Default {
	d := &Default{}
	d.Data.Title = "Default Page"
	d.Data.Name = "default"
	d.Data.Dir = ""

	d.Data.CSS = []render.IRender{render.NewCSS("default.css", "")}

	d.Data.Head = []render.IRender{}
	d.Data.Top = []render.IRender{render.NewTemplate("_header.htm")}

	// requirejs
	data := make(map[string]string, 0)
	data["defer"] = ""
	data["async"] = "true"

	data["data-main"] = httpd.Site.AssetsURL + util.PathToURL("/assets/js/", httpd.RunMode+"/page", d.Dir(), d.Name())
	data["src"] = "/assets/js/require.js"
	d.Data.Bottom = []render.IRender{render.NewTemplate("_bottom.htm"), render.Script("", data)}

	return d
}

// Page page
func (d *Default) Page(ctx *fasthttp.RequestCtx) error {
	fmt.Println("rounter params:", ctx.UserValue("name"), ctx.Value("name"))
	return nil
}

// Auth page
func (d *Default) Auth(ctx *fasthttp.RequestCtx) (string, error) {
	return "", nil
}
