package page

import "github.com/kere/gno"

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
	return d
}

// Prepare page
func (d *Default) Prepare() error {
	d.AddHead("<meta charset=\"utf-8\">")
	d.AddJS("default.js")
	d.AddCSS("default.css")

	d.AddTop("_header.htm", nil)
	d.AddBottom("_bottom.htm", nil)

	d.Data = DefaultData{Name: "tom"}
	return nil
}

// DefaultData is page data
type DefaultData struct {
	Name string
}
