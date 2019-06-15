package websock

import (
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

// WS class
type WS struct {
}

// NewWS f
func NewWS() *WS {
	return &WS{}
}

// Auth f
func (w *WS) Auth(ctx *fasthttp.RequestCtx) error {
	return nil
}

// SayHi f
func (w *WS) SayHi(args util.MapData) (interface{}, error) {
	return "ni hao a...", nil
}
