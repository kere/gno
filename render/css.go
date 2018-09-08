package render

import (
	"io"
	"os"
	"strings"
)

// CSS class
type CSS struct {
	Theme    string
	FileName string
	Data     map[string]string
}

// NewCSS new
func NewCSS(fileName, theme string) CSS {
	return CSS{FileName: fileName, Theme: theme}
}

// Render func
func (t CSS) Render(w io.Writer) error {
	w.Write(bCSSTagBegin)

	filename := t.FileName
	if strings.HasPrefix(filename, "http") {
		w.Write([]byte(filename))
		if CSSVersion != "" {
			w.Write([]byte("?v=" + CSSVersion))
		}
		w.Write(bCSSTagEnd)
		return nil
	}

	if os.PathSeparator == '\\' {
		filename = strings.Replace(t.FileName, "\\", "/", -1)
	}

	if t.Theme == "" {
		w.Write([]byte(AssetsURL + "/assets/css/" + filename))
	} else {
		w.Write([]byte(AssetsURL + "/assets/css/" + t.Theme + "/" + filename))
	}

	if CSSVersion != "" {
		w.Write([]byte("?v=" + CSSVersion))
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
