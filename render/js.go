package render

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// JS class
type JS struct {
	FileName string
	Data     map[string]string
}

// NewJS new
func NewJS(fileName string) *JS {
	return &JS{FileName: fileName}
}

// Render f
func (t *JS) Render(w io.Writer) error {
	w.Write(bJsTagBegin)

	filename := t.FileName
	if os.PathSeparator == '\\' {
		filename = strings.Replace(t.FileName, "\\", "/", -1)
	}

	w.Write([]byte(AssetsURL + "/assets/js/" + RunMode + "/" + filename))

	if JSVersion != "" {
		w.Write([]byte(fmt.Sprint("?v=", JSVersion)))
	}

	w.Write(BytesQuote)
	if t.Data != nil {
		for k, val := range t.Data {
			w.Write(BytesSpace)
			w.Write([]byte(k))
			w.Write(BytesEqual)
			w.Write(BytesQuote)
			w.Write([]byte(val))
			w.Write(BytesQuote)
		}
	}
	w.Write(bJsTagEnd)
	return nil
}
