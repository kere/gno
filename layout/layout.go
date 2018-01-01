package layout

import (
	"io"

	"github.com/kere/goo/render"
)

var (
	// RunMode dev pro
	RunMode = "dev"
)

var (
	// BytesHTMLBegin bytes
	BytesHTMLBegin = []byte("<!DOCTYPE HTML>\n<html>\n")
	// BytesHTMLEnd bytes
	BytesHTMLEnd = []byte("</html>\n")
	// BytesHTMLBodyBegin bytes
	BytesHTMLBodyBegin = []byte("\n<body>\n")
	// BytesHTMLBodyEnd bytes
	BytesHTMLBodyEnd = []byte("\n</body>\n")

	bRenderS1 = []byte("\n<script>var MYENV='")
	bRenderS2 = []byte("',THEME='")
	bRenderS3 = []byte("'</script>")
)

// Page layout class
type Page struct {
	Title  string
	Theme  string
	Head   *Head
	Top    []render.IRender
	Body   []render.IRender
	Bottom []render.IRender
}

// NewPage func
func NewPage(title string) *Page {
	return &Page{Head: NewHead(title)}
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
	w.Write(BytesHTMLBegin)
	p.Head.Render(w)

	// <body>
	w.Write(BytesHTMLBodyBegin)

	var err error
	for _, r := range p.Top {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	for _, r := range p.Body {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	for _, r := range p.Bottom {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	w.Write(bRenderS1)
	w.Write([]byte(RunMode))
	w.Write(bRenderS2)
	w.Write([]byte(p.Theme))
	w.Write(bRenderS3)
	p.Head.RenderBottomJS(w)

	w.Write(BytesHTMLBodyEnd)
	w.Write(BytesHTMLEnd)

	return nil
}
