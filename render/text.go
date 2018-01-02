package render

import (
	"html/template"
	"io"
)

// Text class
type Text struct {
	Source string
	Data   map[string]interface{}
}

// NewText new
func NewText(txt string) *Text {
	return &Text{Source: txt}
}

// Render func
func (t *Text) Render(w io.Writer) error {
	tmpl, _ := template.New("").Parse(t.Source)
	return tmpl.Execute(w, t.Data)
}

// String class
type String struct {
	Source string
}

// NewString new
func NewString(src string) *String {
	return &String{Source: src}
}

// Render func
func (t *String) Render(w io.Writer) error {
	w.Write([]byte(t.Source))
	return nil
}
