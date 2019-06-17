package httpd

import (
	"github.com/kere/gno/render"
	"github.com/valyala/fasthttp"
)

// PageData class
type PageData struct {
	Name, Dir     string
	Title         string
	Secret, Nonce string

	Head   []render.IRender
	CSS    []render.IRender
	JS     []render.IRender
	Top    []render.IRender
	Bottom []render.IRender
}

// PageBase class
type PageBase struct {
	Data PageData
}

// Title page
func (d *PageBase) Title() string {
	return d.Data.Title
}

// GetData page
func (d *PageBase) GetData() *PageData {
	return &d.Data
}

// Dir page
func (d *PageBase) Dir() string {
	return d.Data.Dir
}

// Name page
func (d *PageBase) Name() string {
	return d.Data.Name
}

// CSS page
func (d *PageBase) CSS() []render.IRender {
	return d.Data.CSS
}

// JS page
func (d *PageBase) JS() []render.IRender {
	return d.Data.JS
}

// Head page
func (d *PageBase) Head() []render.IRender {
	return d.Data.Head
}

// Top page
func (d *PageBase) Top() []render.IRender {
	return d.Data.Top
}

// Body page
func (d *PageBase) Body() []render.IRender {
	return nil
}

// Bottom page
func (d *PageBase) Bottom() []render.IRender {
	return d.Data.Bottom
}

// CacheOption page
func (d *PageBase) CacheOption() PageCacheOption {
	return PageCacheOption{Mode: CacheModePagePath, StoreMode: CacheStoreMem, HeadExpires: 0, CacheExpires: 0}
}

// Page page
func (d *PageBase) Page(ctx *fasthttp.RequestCtx) error {
	return nil
}

// Auth page
func (d *PageBase) Auth(ctx *fasthttp.RequestCtx) (string, error) {
	return "", nil
}
