package render

import (
	"io"
	"os"
	"strings"

	"github.com/kere/gno/libs/util"
)

var (
	bJsTagBegin = []byte("<script type=\"text/javascript\"")
	bJsSrc      = []byte(" src=\"")
	bJsTagEnd   = []byte("</script>\n")
)

// JS class
type JS struct {
	FileName string
	Src      string
	Data     map[string]string
}

// NewJS new
func NewJS(fileName string) JS {
	return JS{FileName: fileName}
}

// NewScript new
func NewScript(src string) JS {
	return JS{Src: src}
}

// Render f
func (t JS) Render(w io.Writer) error {
	return t.RenderWith(w, Opt{})
}

// RenderWith f
func (t JS) RenderWith(w io.Writer, opt Opt) error {
	w.Write(bJsTagBegin)

	if t.FileName != "" {
		w.Write(bJsSrc)
		if strings.HasPrefix(t.FileName, "http:") || strings.HasPrefix(t.FileName, "https:") {
			w.Write(util.Str2Bytes(t.FileName))

		} else {
			var filename string
			if os.PathSeparator == '\\' {
				filename = strings.Replace(t.FileName, "\\", "/", -1)
			}
			w.Write(util.Str2Bytes(opt.AssetsURL))
			w.Write(util.Str2Bytes("/js/"))
			w.Write(util.Str2Bytes(filename))
		}

		if opt.JSVersion != "" {
			w.Write(bVerStr)
			w.Write(util.Str2Bytes(opt.JSVersion))
		}
		w.Write(BytesQuote)
	}

	// write property data
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
	w.Write(BytesLargeThan) // >

	if t.Src != "" {
		w.Write(util.Str2Bytes(t.Src))
	}

	w.Write(bJsTagEnd)

	return nil
}
