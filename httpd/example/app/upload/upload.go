package upload

import "github.com/valyala/fasthttp"

// Image class
type Image struct {
}

func NewImage() *Image {
	return &Image{}
}

// Auth a
func (m *Image) Auth(ctx *fasthttp.RequestCtx) error {
	return nil
}

// Do it
func (m *Image) Do(ctx *fasthttp.RequestCtx) error {
	return nil
}

// StoreDir f
func (m *Image) StoreDir(last []byte) string {
	return "Z:/"
}

// FileName f
func (m *Image) FileName() string {
	return "a.jpg"
}
