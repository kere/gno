package httpd

import (
	"io"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
)

// StringRender class
type StringRender struct {
	Src []byte
}

// NewStrRender new
func NewStrRender(src string) *StringRender {
	return &StringRender{Src: util.Str2Bytes(src)}
}

// NewStrbRender new
func NewStrbRender(src []byte) *StringRender {
	return &StringRender{Src: src}
}

// Render func
func (t *StringRender) Render(w io.Writer) error {
	_, err := w.Write(t.Src)
	return err
}

// // RenderA func
// func (t *StringRender) RenderA(w io.Writer, pa *PageAttr) error {
// 	return t.Render(w)
// }

// BufferRender class
type BufferRender struct {
	Buf *bytebufferpool.ByteBuffer
}

// Release func
func (t *BufferRender) Release() {
	bytebufferpool.Put(t.Buf)
}
