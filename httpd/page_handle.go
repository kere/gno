package httpd

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/kere/gno/libs/log"
	"github.com/valyala/fasthttp"
)

// pageHandle page http handle
func pageHandle(p IPage, ctx *fasthttp.RequestCtx) {
	err := p.Auth(ctx)
	if err != nil {
		u, _ := url.Parse(Site.LoginURL)
		u.Query().Add("url", string(ctx.RequestURI()))
		ctx.Redirect(u.String(), http.StatusSeeOther)
		return
	}

	if TryGetCache(ctx, p) {
		return
	}

	err = p.Page(ctx)
	if err != nil {
		log.App.Error(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	buf := bytes.NewBuffer(nil)

	err = renderPage(buf, p, ctx.URI().PathOriginal())
	if err != nil {
		log.App.Error(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	TrySetCache(ctx, p, buf)
	ctx.Write(buf.Bytes())
	ctx.SetContentTypeBytes(contentTypePage)
}
