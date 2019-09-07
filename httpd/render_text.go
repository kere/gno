package httpd

import (
	"html/template"
	"io"
)

// TextRender class
type TextRender struct {
	Src  string
	Data map[string]interface{}
}

// NewTextRender new
func NewTextRender(txt string) *TextRender {
	return &TextRender{Src: txt}
}

// Render func
func (t *TextRender) Render(w io.Writer) error {
	tmpl, _ := template.New("").Parse(t.Src)
	return tmpl.Execute(w, t.Data)
}
