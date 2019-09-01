package render

import (
	"io"

	"github.com/kere/gno/libs/util"
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
	w.Write(util.Str2Bytes(h.Value))
	return nil
}
