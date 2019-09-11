package httpd

import (
	"io"
	"os"
	"strings"

	"github.com/kere/gno/libs/util"
)

// CSS class
type CSS struct {
	Theme    string
	FileName string
	Data     map[string]string
}

// NewCSS new
func NewCSS(fileName string) *CSS {
	return &CSS{FileName: fileName}
}

// RenderA func
func (t *CSS) RenderA(w io.Writer, pd *PageAttr) error {
	w.Write(bCSSTagBegin)

	filename := t.FileName
	if strings.HasPrefix(filename, "http") {
		w.Write([]byte(filename))
		w.Write(bCSSTagEnd)
		return nil
	}

	if os.PathSeparator == '\\' {
		filename = strings.Replace(t.FileName, "\\", "/", -1)
	}

	// w.Write([]byte(pd.AssetsURL + "/css/" + t.Theme + "/" + filename))
	w.Write(util.Str2Bytes(pd.SiteData.AssetsURL))
	w.Write(util.Str2Bytes("/css/"))
	w.Write(util.Str2Bytes(filename))

	if pd.SiteData.CSSVersion != "" {
		w.Write(bVerStr)
		w.Write(util.Str2Bytes(pd.SiteData.CSSVersion))
	}

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
	w.Write(bCSSTagEnd)
	return nil
}
