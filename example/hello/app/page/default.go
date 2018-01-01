package page

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/goo/page"
)

// Default page class
type Default struct {
	page.Page
}

// NewDefaultPage func
func NewDefaultPage() page.IPage {
	p := &Default{}
	return p
}

// Init page
func (d *Default) Init(req *http.Request, ps httprouter.Params) {
	d.Page.Init(req, ps)
	d.Title = "Default Page"
	d.Name = "default"
	d.Theme = ""
	d.AddHead("<meta charset=\"utf-8\">")
	d.AddJS("default.js")
	d.AddCSS("default.css")

	d.AddTop("_header.htm", nil)
	d.AddBottom("_bottom.htm", nil)
}

// Prepare page
func (d *Default) Prepare() error {
	return nil
}
