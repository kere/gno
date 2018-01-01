package render

import (
	"io"
)

var (
	// RunMode dev pro
	RunMode = ""
	// AssetsURL url
	AssetsURL = ""
	// JSVersion js ?v=001
	JSVersion = ""
	// CSSVersion css ?v=001
	CSSVersion = ""
	// TemplateLeftDelim for template
	TemplateLeftDelim = ""
	// TemplateRightDelim for template
	TemplateRightDelim = ""
)
var (
	// BytesEqual equal
	BytesEqual = []byte("=")
	// BytesQuote quote
	BytesQuote = []byte("\"")
	//BytesSpace space
	BytesSpace = []byte(" ")
	// BytesBreak break
	BytesBreak = []byte("\n")

	bJsTagBegin  = []byte("<script src=\"")
	bJsTagEnd    = []byte("></script>\n")
	bCSSTagBegin = []byte("<link href=\"")
	bCSSTagEnd   = []byte("\" rel=\"stylesheet\"/>\n")
)

// IRender interface
type IRender interface {
	Render(io.Writer) error
}
