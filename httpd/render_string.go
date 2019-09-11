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

// BufferRender class
type BufferRender struct {
	Buf *bytebufferpool.ByteBuffer
}

// Release func
func (t *BufferRender) Release() {
	bytebufferpool.Put(t.Buf)
}
