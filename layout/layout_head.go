package layout

import (
	"io"

	"github.com/kere/gno/render"
)

var (
	bHeadBegin  = []byte("<head>\n")
	bHeadEnd    = []byte("</head>\n")
	bTitleBegin = []byte("<title>")
	bTitleEnd   = []byte("</title>\n")
	metaCharset = []byte("<meta charset=\"UTF-8\">\n")
)

// Head for page
type Head struct {
	JSPosition string
	Title      string
	Theme      string
	HeadItems  []render.IRender
	CSSRenders []render.IRender
	JSRenders  []render.IRender
}

// NewHead new
func NewHead() *Head {
	return &Head{}
}

// AddHeadItem func
func (h *Head) AddHeadItem(r render.IRender) {
	h.HeadItems = append(h.HeadItems, r)
}

// AddJS func
func (h *Head) AddJS(r render.IRender) {
	h.JSRenders = append(h.JSRenders, r)
}

// AddCSS func
func (h *Head) AddCSS(r render.IRender) {
	h.CSSRenders = append(h.CSSRenders, r)
}

// Render func
func (h *Head) Render(w io.Writer) error {
	w.Write(bHeadBegin)
	w.Write(metaCharset)
	var err error

	w.Write(bRenderS1)
	w.Write([]byte(RunMode))
	w.Write(bRenderS2)
	w.Write([]byte(h.Theme))
	w.Write(bRenderS3)

	for _, r := range h.HeadItems {
		if err = r.Render(w); err != nil {
			return err
		}
	}
	for _, r := range h.CSSRenders {
		if err = r.Render(w); err != nil {
			return err
		}
	}
	if h.JSPosition == "head" {
		for _, r := range h.JSRenders {
			if err = r.Render(w); err != nil {
				return err
			}
		}
	}

	w.Write(bTitleBegin)
	w.Write([]byte(h.Title))
	w.Write(bTitleEnd)

	w.Write(bHeadEnd)
	return nil
}

// RenderBottomJS func
func (h *Head) RenderBottomJS(w io.Writer) error {
	if h.JSPosition == "head" {
		return nil
	}

	var err error
	for _, r := range h.JSRenders {
		if err = r.Render(w); err != nil {
			return err
		}
	}
	return nil
}
