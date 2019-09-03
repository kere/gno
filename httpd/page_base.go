package httpd

import (
	"github.com/kere/gno/httpd/render"
	"github.com/valyala/fasthttp"
)

const (
	// JSPositionHead top head
	JSPositionHead = 1
	// JSPositionBottom bottom
	JSPositionBottom = 9
)

// PageData class
type PageData struct {
	Name, Dir string
	Title     string
	// Secret, Nonce string
	CacheOption PageCacheOption
	JSPosition  int

	Head []render.IRender

	JS  []render.IRenderWith
	CSS []render.IRenderWith

	Top    []render.IRender
	Body   []render.IRender
	Bottom []render.IRender
}

// Init page
func (d *PageData) Init(title, name, dir string) {
	d.Title = title
	d.Name = name
	d.Dir = dir
}

// P class
type P struct {
	D PageData
}

// Data page
func (d *P) Data() *PageData {
	return &d.D
}

// SetData page
func (d *P) SetData(pd *PageData) {
	d.D = *pd
}

// Page do
func (d *P) Page(ctx *fasthttp.RequestCtx) error {
	return nil
}

// Auth page
func (d *P) Auth(ctx *fasthttp.RequestCtx) error {
	return nil
}
