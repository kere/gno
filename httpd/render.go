package httpd

import (
	"io"
)

var (
	// // JSVersion js ?v=001
	// JSVersion []byte
	//
	// // CSSVersion js ?v=001
	// CSSVersion []byte
	//
	// // AssetsURL url
	// AssetsURL = ""

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

	// BJsTagBegin bytes
	BJsTagBegin = []byte("<script type=\"text/javascript\"")
	// BJsSrc bytes
	BJsSrc = []byte(" src=\"")
	// BJsTagEnd bytes
	BJsTagEnd = []byte("</script>\n")
)

// IRender interface
type IRender interface {
	Render(io.Writer) error
}

// IRenderA with Attr
type IRenderA interface {
	RenderA(io.Writer, *PageAttr) error
}

// IRenderD with data
type IRenderD interface {
	RenderD(io.Writer, interface{}) error
}

// // IRenderAD with Attr and Data
// type IRenderAD interface {
// 	RenderD(io.Writer, *PageAttr, *PageData) error
// }
