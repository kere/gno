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
