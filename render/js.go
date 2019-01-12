package render

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	bJsTagBegin = []byte("<script text=\"text/javascript\" src=\"")
	bJsTagEnd   = []byte("></script>\n")
)

// JS class
type JS struct {
	FileName string
	Data     map[string]string
}

// NewJS new
func NewJS(fileName string) JS {
	return JS{FileName: fileName}
}

// Render f
func (t JS) Render(w io.Writer) error {
	w.Write(bJsTagBegin)

	filename := t.FileName
	if strings.HasPrefix(filename, "http:") || strings.HasPrefix(filename, "https:") {
		w.Write([]byte(filename))

	} else {
		if os.PathSeparator == '\\' {
			filename = strings.Replace(t.FileName, "\\", "/", -1)
		}

		w.Write([]byte(AssetsURL + "/assets/js/" + filename))
	}
	if JSVersion != "" {
		w.Write([]byte(fmt.Sprint("?gv=", JSVersion)))
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
