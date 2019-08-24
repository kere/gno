package httpd

import (
	"github.com/kere/gno/httpd/render"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

// PageData class
type PageData struct {
	Name, Dir     string
	Title         []byte
	Secret, Nonce string
	CacheOption   PageCacheOption

	Head   []render.IRender
	CSS    []render.IRender
	JS     []render.IRender
	Top    []render.IRender
	Body   []render.IRender
	Bottom []render.IRender
}

// Init page
func (d *PageData) Init(title, name, dir string) {
	d.Title = util.Str2Bytes(title)
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
