package httpd

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/render"
	"github.com/valyala/fasthttp"
)

// IPage interface
type IPage interface {
	Title() string
	Lang() string

	Name() string
	Dir() string

	Head() []render.IRender
	JS() []render.IRender
	CSS() []render.IRender

	Top() []render.IRender
	Body() []render.IRender
	Bottom() []render.IRender

	CacheOption() PageCacheOption

	Auth(ctx *fasthttp.RequestCtx) (string, error)
	// Before(ctx *fasthttp.RequestCtx) error
	Page(ctx *fasthttp.RequestCtx) error
}

// RegistGet router
func (s *SiteServer) RegistGet(rule string, p IPage) {
	s.Router.GET(rule, func(ctx *fasthttp.RequestCtx) {
		loginURL, err := p.Auth(ctx)
		if loginURL != "" && err != nil {
			u, _ := url.Parse(loginURL)
			u.Query().Add("url", string(ctx.RequestURI()))
			ctx.Redirect(u.String(), http.StatusSeeOther)
			return
		}

		if TryGetCache(ctx, p) {
			log.App.Debug("Page Cache", string(ctx.URI().Path()))
			return
		}

		err = p.Page(ctx)
		if err != nil {
			doPageErr(s.ErrorURL, ctx, err)
			return
		}

		buf := bytes.NewBuffer(nil)

		err = renderPage(s, buf, p, ctx.URI().PathOriginal())

		if err != nil {
			doPageErr(s.ErrorURL, ctx, err)
			return
		}

		TrySetCache(ctx, p, buf)
		ctx.Write(buf.Bytes())
		ctx.SetContentTypeBytes(contentTypePage)
		// _, err = p.GetResponseWriter().Write(buf.Bytes())
	})
}

// RegistPost router
func (s *SiteServer) RegistPost(rule string, p IPage) {
	s.Router.POST(rule, func(ctx *fasthttp.RequestCtx) {
	})
}

func doPageErr(errorURL string, ctx *fasthttp.RequestCtx, err error) {
	if errorURL == "" {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	// ErrorURL redirect to
	ctx.Redirect(errorURL+"?msg="+err.Error(), http.StatusSeeOther)
}

// // SetPageToken token
// func (s *SiteServer) SetPageToken(ctx *fasthttp.RequestCtx, p IPage) {
// 	token := s.PageToken(ctx.URI().Path(), fmt.Sprint(time.Now().Unix()))
//
// 	token = base64.StdEncoding.EncodeToString([]byte(token))
// 	ctx.Request.Header.SetCookie(PageAccessTokenField, url.PathEscape(token))
// }
