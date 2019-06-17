package httpd

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

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

	println(1)
	if TryCache(ctx, p) {
		fmt.Println("Header:", ctx.Request.Header.String())
		return
	}
	println(2)

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
	if p.Data().CacheOption.StoreMode == CacheStoreNone {
		ctx.Request.Header.Set(HeaderCacheCtl, headValCacheNone)
	} else {
		ctx.Request.Header.Set(HeaderCacheCtl, fmt.Sprint(headValCache, p.Data().CacheOption.HeadExpires))
		ctx.Request.Header.Set(HeaderLastModified, lastModified)
	}
}
