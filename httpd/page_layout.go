package httpd

import (
	"io"
	"path/filepath"

	"github.com/kere/gno/render"
)

var (
	bHeadBegin  = []byte("<head>\n")
	bHeadEnd    = []byte("</head>\n")
	bTitleBegin = []byte("<title>")
	bTitleEnd   = []byte("</title>\n")

	metaCharset     = []byte("<meta charset=\"UTF-8\">\n")
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
func renderPage(site *SiteServer, w io.Writer, p IPage, bPath []byte) error {
	// <html>
	w.Write(bytesHTMLBegin)
	lang := p.Lang()
	if lang == "" {
		w.Write([]byte("en"))
	} else {
		w.Write([]byte(p.Lang()))
	}
	w.Write(bytesHTMLBegin2)

	// head -------------------------
	w.Write(bHeadBegin)
	w.Write(metaCharset)

	w.Write(bTitleBegin)
	w.Write([]byte(p.Title()))
	w.Write(bTitleEnd)

	w.Write(bRenderS1)
	w.Write([]byte(RunMode))
	w.Write(bRenderS2)

	token := buildToken(bPath, site.Secret, site.Salt)

	w.Write([]byte(token))

	w.Write(bRenderS3)
	heads := p.Head()
	for _, r := range heads {
		if err := r.Render(w); err != nil {
			return err
		}
	}
	css := p.CSS()
	for _, r := range css {
		if err := r.Render(w); err != nil {
			return err
		}
	}
	js := p.JS()
	for _, r := range js {
		if err := r.Render(w); err != nil {
			return err
		}
	}
	w.Write(bHeadEnd)

	// <body>
	w.Write(bytesHTMLBodyBegin)

	var err error
	top := p.Top()
	for _, r := range top {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	body := p.Body()
	if len(body) == 0 {
		r := render.NewTemplate(filepath.Join(p.Dir(), p.Name()+defaultTemplateSubfix))
		if err = r.Render(w); err != nil {
			return err
		}
	} else {
		for _, r := range body {
			if err = r.Render(w); err != nil {
				return err
			}
		}
	}

	bottom := p.Bottom()
	for _, r := range bottom {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	w.Write(bytesHTMLBodyEnd)
	w.Write(bytesHTMLEnd)

	return nil
}
