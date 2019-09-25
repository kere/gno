package httpd

import (
	"io"

	"github.com/kere/gno/libs/util"
)

var (
	bHeadBegin  = []byte("<head>\n")
	bHeadEnd    = []byte("</head>\n")
	bTitleBegin = []byte("<title>")
	bTitleEnd   = []byte("</title>\n")

	// <meta http-equiv="content-type" content="txt/html; charset=utf-8"/>
	metaCharset     = []byte("<meta http-equiv=\"content-type\" content=\"txt/html; charset=utf-8\"/>\n")
	bytesHTMLBegin  = []byte("<!DOCTYPE HTML>\n<html lang=\"")
	bytesHTMLBegin2 = []byte("\">\n")
	bytesHTMLEnd    = []byte("</html>\n")
	// BytesHTMLBodyBegin bytes
	bytesHTMLBodyBegin = []byte("\n<body>\n")
	// BytesHTMLBodyEnd bytes
	bytesHTMLBodyEnd = []byte("\n</body>\n")

	bRenderS1 = []byte("\n<script type=\"text/javascript\">var MYENV='")
	bRenderS2 = []byte("'," + PageAccessTokenField + "='")
	bRenderS3 = []byte("';</script>")

	contentTypePage = []byte("text/html; charset=utf-8")
)

// renderPage func
func renderPage(w io.Writer, pa *PageAttr, dat interface{}, bPath []byte) error {
	// <html>
	w.Write(bytesHTMLBegin)
	w.Write(util.Str2Bytes(pa.SiteData.Lang))
	w.Write(bytesHTMLBegin2)

	// head -------------------------
	w.Write(bHeadBegin)
	w.Write(metaCharset)

	w.Write(bTitleBegin)
	// if pageData.Title != "" {
	// 	w.Write(util.Str2Bytes(pageData.Title))
	// } else {
	w.Write(util.Str2Bytes(pa.Title))
	// }

	w.Write(bTitleEnd)

	w.Write(bRenderS1)
	w.Write([]byte(RunMode))
	w.Write(bRenderS2)

	token := buildToken(bPath, pa.SiteData.Secret, pa.SiteData.Nonce)
	w.Write(util.Str2Bytes(token))

	w.Write(bRenderS3)
	// Head
	for _, r := range pa.Head {
		if err := r.Render(w); err != nil {
			return err
		}
	}

	// CSS
	for _, r := range pa.CSS {
		if err := r.RenderA(w, pa); err != nil {
			return err
		}
	}
	// JS
	if pa.JSPosition == JSPositionHead {
		for _, r := range pa.JS {
			if err := r.RenderA(w, pa); err != nil {
				return err
			}
		}
	}

	w.Write(bHeadEnd)
	w.Write(bytesHTMLBodyBegin) // <body>

	// Top
	var err error
	for _, r := range pa.Top {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	// Body
	if pa.Body != nil {
		if err = pa.Body.RenderD(w, dat); err != nil {
			return err
		}
	}

	// Bottom JS
	if pa.JSPosition == JSPositionBottom {
		for _, r := range pa.JS {
			if err := r.RenderA(w, pa); err != nil {
				return err
			}
		}
	}

	// Bottom
	for _, r := range pa.Bottom {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	w.Write(bytesHTMLBodyEnd)
	w.Write(bytesHTMLEnd)

	return nil
}
