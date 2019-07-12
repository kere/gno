package render

import (
	"io"
)

var (
	// JSVersion js ?v=001
	JSVersion []byte

	// CSSVersion js ?v=001
	CSSVersion []byte

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
	// BytesLargeThan >
	BytesLargeThan = []byte(">")
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