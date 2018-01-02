package page

import (
	"github.com/kere/goo"
)

// Default page class
type Default struct {
	goo.Page
}

// NewDefaultPage func
func NewDefaultPage() goo.IPage {
	p := &Default{}
	return p
}

// Build page
func (d *Default) Build() error {
	d.Title = "Default Page"
	d.Name = "default"
	d.Theme = ""
	d.AddHead("<meta charset=\"utf-8\">")
	d.AddJS("default.js")
	d.AddCSS("default.css")

	d.AddTop("_header.htm", nil)
	d.AddBottom("_bottom.htm", nil)

	return nil
}

// Prepare page
func (d *Default) Prepare() error {
	return nil
}
