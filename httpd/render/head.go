package render

import (
	"io"
)

// Head render
type Head struct {
	Value string
}

// NewHead render
func NewHead(v string) Head {
	return Head{Value: v}
}

// Render func
func (h Head) Render(w io.Writer) error {
	w.Write([]byte(h.Value))
	return nil
}
