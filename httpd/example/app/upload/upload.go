package upload

import "github.com/valyala/fasthttp"

// Image class
type Image struct {
}

// NewImage img
func NewImage() *Image {
	return &Image{}
}

// Auth a
func (m *Image) Auth(ctx *fasthttp.RequestCtx) error {
	return nil
}

// Success it
func (m *Image) Success(name, folder string) error {
	return nil
}

// StoreDir f
func (m *Image) StoreDir() string {
	return "Z:/"
}

// FileName f
func (m *Image) FileName() string {
	return "a.jpg"
}
