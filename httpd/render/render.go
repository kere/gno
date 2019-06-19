package render

import (
	"io"
)

// JSVersion js ?v=001
var JSVersion []byte

// CSSVersion js ?v=001
var CSSVersion []byte

var (
	// AssetsURL url
	AssetsURL = ""

	// TemplateLeftDelim for template
	TemplateLeftDelim = ""
	// TemplateRightDelim for template
	TemplateRightDelim = ""
	// BytesEqual equal
	BytesEqual = []byte("=")
	// BytesQuote quote
	BytesQuote = []byte("\"")
	//BytesSpace space
	BytesSpace = []byte(" ")
	// BytesBreak break
	BytesBreak = []byte("\n")

	bCSSTagBegin = []byte("<link href=\"")
	bCSSTagEnd   = []byte("\" rel=\"stylesheet\"/>\n")
	bVerStr      = []byte("?gv=")
)

// IRender interface
type IRender interface {
	Render(io.Writer) error
}
