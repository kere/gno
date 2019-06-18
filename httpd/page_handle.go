package httpd

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

// pageHandle page http handle
func pageHandle(p IPage, ctx *fasthttp.RequestCtx) {
	err := p.Auth(ctx)
	if err != nil {
		u, _ := url.Parse(Site.LoginURL)
		u.Query().Add("url", string(ctx.RequestURI()))
		ctx.Redirect(u.String(), fasthttp.StatusSeeOther)
		return
	}

	println(1)
	if TryCache(ctx, p) {
		fmt.Println("Header:", ctx.Request.Header.String())
		return
	}
	println(2)

	err = p.Page(ctx)
	if err != nil {
		log.App.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	buf := bytes.NewBuffer(nil)

	err = renderPage(buf, p, ctx.URI().PathOriginal())
	if err != nil {
		log.App.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	TrySetCache(ctx, p, buf)

	ctx.Write(buf.Bytes())
}

func gmtNowTime(d time.Time) string {
	lc, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
	gmt := d.In(lc)

	return gmt.Format(LastModifiedFormat)
}

func setHeader(p IPage, ctx *fasthttp.RequestCtx, lastModified string) {
	ctx.SetContentTypeBytes(contentTypePage)
	// set response header
	// Cache-Control: public, max-age=3600
	// must-revalidate

	mode := p.Data().CacheOption.HTTPHead
	switch {
	case mode == 1:
		ctx.Response.Header.SetBytesK(HeaderEtag, fmt.Sprintf("%x", util.MD5(lastModified)))
		ctx.Response.Header.SetBytesK(HeaderLastModified, lastModified)
	case mode > 1:
		ctx.Response.Header.SetBytesK(HeaderCacheCtl, fmt.Sprint(headSValMaxAge, mode))
		ctx.Response.Header.SetBytesK(HeaderLastModified, lastModified)
	default:
		ctx.Response.Header.SetBytesK(HeaderCacheCtl, headSValNoCache)
	}
}
