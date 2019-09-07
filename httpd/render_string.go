package httpd

import (
	"io"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
)

// StringRender class
type StringRender struct {
	Src string
}

// NewStrRender new
func NewStrRender(src string) *StringRender {
	return &StringRender{Src: src}
}

// Render func
func (t *StringRender) Render(w io.Writer) error {
	w.Write(util.Str2Bytes(t.Src))
	return nil
}

// RenderWith func
func (t *StringRender) RenderWith(w io.Writer, pd *PageData) error {
	w.Write(util.Str2Bytes(t.Src))
	return nil
}

// BufferRender class
type BufferRender struct {
	Buf *bytebufferpool.ByteBuffer
}

// // Render func
// func (t *BufferRender) Render(w io.Writer) error {
// 	w.Write(t.Buf.Bytes())
// 	return nil
// }
//
// // RenderWith func
// func (t *BufferRender) RenderWith(w io.Writer, pd *PageData) error {
// 	w.Write(t.Buf.Bytes())
// 	return nil
// }

// Release func
func (t *BufferRender) Release() {
	bytebufferpool.Put(t.Buf)
}
