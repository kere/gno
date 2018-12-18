package page

import (
	"github.com/kere/gno"
	"github.com/kere/gno/libs/util"
)

// Default page class
type Default struct {
	gno.Page
}

// NewDefaultPage func
func NewDefaultPage() gno.IPage {
	d := &Default{}
	d.Title = "Default Page"
	d.Name = "default"
	d.Dir = ""
	d.Theme = ""
	// d.SetPageExpires(gno.CacheModePage, 300)
	return d
}

// Prepare page
func (d *Default) Prepare() error {
	d.AddHead("<meta charset=\"utf-8\">")
	d.AddCSS("default.css")

	d.AddTop("_header.htm", nil)
	d.AddBottom("_bottom.htm", nil)

	d.Data = DefaultData{Name: "tom"}

	// requirejs
	data := make(map[string]string, 0)
	data["defer"] = ""
	data["async"] = "true"

	data["data-main"] = gno.Site.AssetsURL + util.PathToURL("/assets/js/", gno.RunMode+"/page", d.GetDir(), d.GetName())
	data["src"] = "/assets/js/require.js"
	d.AddScript("bottom", "", data)
	return nil
}

// DefaultData is page data
type DefaultData struct {
	Name string
}
