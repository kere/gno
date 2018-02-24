package layout

import (
	"io"

	"github.com/kere/gno/render"
)

var (
	// RunMode dev pro
	RunMode = "dev"
)

var (
	bytesHTMLBegin  = []byte("<!DOCTYPE HTML>\n<html lang=\"")
	bytesHTMLBegin2 = []byte("\">\n")
	bytesHTMLEnd    = []byte("</html>\n")
	// BytesHTMLBodyBegin bytes
	bytesHTMLBodyBegin = []byte("\n<body>\n")
	// BytesHTMLBodyEnd bytes
	bytesHTMLBodyEnd = []byte("\n</body>\n")

	bRenderS1 = []byte("\n<script type=\"text/javascript\">var MYENV='")
	bRenderS2 = []byte("',THEME='")
	bRenderS3 = []byte("';</script>")
)

// Page layout class
type Page struct {
	Lang   string
	Head   *Head
	Top    []render.IRender
	Body   []render.IRender
	Bottom []render.IRender
}

// NewPage func
func NewPage() *Page {
	return &Page{Head: NewHead()}
}

// AddTop add top render template
func (p *Page) AddTop(filename string, data interface{}) {
	r := render.NewTemplate(filename)
	r.Data = data
	p.Top = append(p.Top, r)
}

// AddBody add body render template
func (p *Page) AddBody(filename string, data interface{}) {
	r := render.NewTemplate(filename)
	r.Data = data
	p.Body = append(p.Body, r)
}

// AddBottom add bottom render template
func (p *Page) AddBottom(filename string, data interface{}) {
	r := render.NewTemplate(filename)
	r.Data = data
	p.Bottom = append(p.Bottom, r)
}

// Render func
func (p *Page) Render(w io.Writer) error {
	// <html>
	w.Write(bytesHTMLBegin)
	w.Write([]byte(p.Lang))
	w.Write(bytesHTMLBegin2)
	p.Head.Render(w)

	// <body>
	w.Write(bytesHTMLBodyBegin)

	var err error
	for _, r := range p.Top {
		if err = r.Render(w); err != nil {
			// println("top")
			return err
		}
	}

	for _, r := range p.Body {
		if err = r.Render(w); err != nil {
			// println("body")
			return err
		}
	}

	for _, r := range p.Bottom {
		if err = r.Render(w); err != nil {
			// println("bottom")
			return err
		}
	}

	p.Head.RenderBottomJS(w)

	w.Write(bytesHTMLBodyEnd)
	w.Write(bytesHTMLEnd)

	return nil
}
