package httpd

import (
	"io"
	"os"
	"strings"

	"github.com/kere/gno/libs/util"
)

// JS class
type JS struct {
	FileName string
	Src      []byte
	SrcList  [][]byte
	Attr     [][2]string
}

// NewJS new
func NewJS(fileName string) *JS {
	return &JS{FileName: fileName}
}

// NewJSSrc new
func NewJSSrc(src string, attr [][2]string) *JS {
	return &JS{Src: util.Str2Bytes(src), Attr: attr}
}

// NewJSSrcB new
func NewJSSrcB(src []byte, attr [][2]string) *JS {
	return &JS{Src: src, Attr: attr}
}

// Render js
func (t *JS) Render(w io.Writer) error {
	return t.RenderA(w, nil)
}

// RenderA with page attr
func (t *JS) RenderA(w io.Writer, pd *PageAttr) error {
	w.Write(BJsTagBegin)

	if t.FileName != "" {
		w.Write(BJsSrc)
		if strings.HasPrefix(t.FileName, "http:") || strings.HasPrefix(t.FileName, "https:") {
			w.Write(util.Str2Bytes(t.FileName))

		} else {
			filename := t.FileName
			if os.PathSeparator == '\\' {
				filename = strings.Replace(t.FileName, "\\", "/", -1)
			}
			w.Write(util.Str2Bytes(pd.SiteData.AssetsURL))
			w.Write(util.Str2Bytes("/js/"))
			w.Write(util.Str2Bytes(filename))
			if pd.SiteData.JSVersion != "" {
				w.Write(bVerStr)
				w.Write(util.Str2Bytes(pd.SiteData.JSVersion))
			}
		}

		w.Write(BytesQuote)
	}

	// write property data
	if len(t.Attr) > 0 {
		for _, val := range t.Attr {
			w.Write(BytesSpace)
			w.Write(util.Str2Bytes(val[0]))
			w.Write(BytesEqual)
			w.Write(BytesQuote)
			w.Write(util.Str2Bytes(val[1]))
			w.Write(BytesQuote)
		}
	}
	w.Write(BytesLargeThan) // >

	if len(t.Src) > 0 {
		w.Write(t.Src)
	}

	w.Write(BJsTagEnd)

	return nil
}

// RequireJS class
type RequireJS struct {
	JS
}

var requireJSArg []byte

// RenderA with page attr
func (t *RequireJS) RenderA(w io.Writer, pd *PageAttr) error {
	w.Write(BJsTagBegin)

	if t.FileName != "" {
		w.Write(BJsSrc)
		if strings.HasPrefix(t.FileName, "http:") || strings.HasPrefix(t.FileName, "https:") {
			w.Write(util.Str2Bytes(t.FileName))

		} else {
			filename := t.FileName
			if os.PathSeparator == '\\' {
				filename = strings.Replace(t.FileName, "\\", "/", -1)
			}
			w.Write(util.Str2Bytes(pd.SiteData.AssetsURL))
			w.Write(util.Str2Bytes("/js/"))
			w.Write(util.Str2Bytes(filename))
			if pd.SiteData.JSVersion != "" {
				w.Write(bVerStr)
				w.Write(util.Str2Bytes(pd.SiteData.JSVersion))
			}
		}

		w.Write(BytesQuote)
	}

	// write property data
	if len(t.Attr) > 0 {
		for _, val := range t.Attr {
			w.Write(BytesSpace)
			w.Write(util.Str2Bytes(val[0]))
			w.Write(BytesEqual)
			w.Write(BytesQuote)
			w.Write(util.Str2Bytes(val[1]))
			w.Write(BytesQuote)
		}
	}
	w.Write(BytesLargeThan) // >

	if len(t.Src) > 0 {
		w.Write(t.Src)
	}

	if pd.SiteData.JSVersion != "" {
		w.Write(util.Str2Bytes(delim2))
		if len(requireJSArg) == 0 {
			requireJSArg = []byte("requirejs.config({urlArgs:'" + pd.SiteData.JSVersion + "'})")
		}
		w.Write(requireJSArg)
	}

	w.Write(BJsTagEnd)

	return nil
}
