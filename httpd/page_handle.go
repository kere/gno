package httpd

import (
	"fmt"
	"net/url"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

// pageHandle page http handle
func pageHandle(site *SiteServer, p IPage, ctx *fasthttp.RequestCtx) {
	err := p.Auth(ctx)
	if err != nil {
		u, _ := url.Parse(site.LoginURL)
		u.Query().Add(sAuthURL, string(ctx.RequestURI()))
		ctx.Redirect(u.String(), fasthttp.StatusSeeOther)
		return
	}

	if TryCache(ctx, p) {
		return
	}

	err = p.Page(ctx)
	if err != nil {
		log.App.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = renderPage(site, ctx, p.Data(), ctx.URI().PathOriginal())
	if err != nil {
		log.App.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = TrySetCache(ctx, p, ctx.Response.Body())
	if err != nil {
		log.App.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func setHeader(p IPage, ctx *fasthttp.RequestCtx, lastModified string) {
	ctx.SetContentTypeBytes(contentTypePage)
	// set response header
	// Cache-Control: public, max-age=3600
	// must-revalidate

	mode := p.Data().CacheOption.HTTPHead
	switch {
	case mode == 1:
		ctx.Response.Header.Set(fasthttp.HeaderETag, fmt.Sprintf("%x", util.MD5(lastModified)))
		ctx.Response.Header.Set(fasthttp.HeaderLastModified, lastModified)
	case mode > 1:
		ctx.Response.Header.Set(fasthttp.HeaderCacheControl, fmt.Sprint(headSValMaxAge, mode))
		ctx.Response.Header.Set(fasthttp.HeaderLastModified, lastModified)
	default:
		ctx.Response.Header.Set(fasthttp.HeaderCacheControl, headSValNoCache)
	}
}
