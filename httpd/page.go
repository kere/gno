package httpd

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/kere/gno/libs/log"
	"github.com/valyala/fasthttp"
)

const (
	// JSPositionHead top head
	JSPositionHead = 1
	// JSPositionBottom bottom
	JSPositionBottom = 9
)

// IPage interface
type IPage interface {
	// Data() *PageData
	Auth(ctx *fasthttp.RequestCtx) error
	Page(ctx *fasthttp.RequestCtx) (interface{}, error)
	Attr() *PageAttr
}

// PageAttr static data
type PageAttr struct {
	Name, Dir string
	Title     string
	// Secret, Nonce string
	CacheOption PageCacheOption
	JSPosition  int

	SiteData *SiteData
	ErrorURL string
	LoginURL string

	Head   []IRender
	JS     []IRenderA
	CSS    []IRenderA
	Body   IRenderD
	Top    []IRender
	Bottom []IRender
}

// P class
type P struct {
	PA PageAttr
}

// Attr page
func (p *P) Attr() *PageAttr {
	return &p.PA
}

// initPage func
func initPage(p IPage) {
	typ := reflect.TypeOf(p)
	s := strings.TrimPrefix(typ.String(), "*")
	arr := strings.Split(s, ".")
	dir := arr[0]
	name := arr[1]

	pa := p.Attr()
	pa.Name = name
	pa.Dir = dir
	if pa.Body == nil {
		pa.Body = NewSiteTemplate(dir, name)
	}
}

// Page do
func (p *P) Page(ctx *fasthttp.RequestCtx) (interface{}, error) {
	return nil, nil
}

// Auth page
func (p *P) Auth(ctx *fasthttp.RequestCtx) error {
	return nil
}

// RegistGet router
func (s *SiteServer) RegistGet(rule string, p IPage) {
	s.Router.GET(rule, func(ctx *fasthttp.RequestCtx) {
		pa := p.Attr()
		pa.SiteData = s.SiteData
		// do auth
		err := p.Auth(ctx)
		if err != nil {
			var loginURL string
			if len(loginURL) == 0 {
				loginURL = pa.SiteData.LoginURL
			} else {
				loginURL = pa.LoginURL
			}
			u, _ := url.Parse(loginURL)
			dat := u.Query()
			dat.Add(sAuthURL, string(ctx.RequestURI()))
			u.RawQuery = dat.Encode()
			ctx.Redirect(u.String(), fasthttp.StatusSeeOther)
			return
		}

		ctx.SetContentTypeBytes(contentTypePage)

		// try cache
		if TryCache(ctx, p) {
			tryStaticHTML(ctx, p)
			return
		}

		// do page
		pdat, err := p.Page(ctx)
		if err != nil {
			doPageErr(pa, ctx, err)
			return
		}

		err = renderPage(ctx, pa, pdat, ctx.URI().PathOriginal())

		if err != nil {
			doPageErr(pa, ctx, err)
			// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		tryStaticHTML(ctx, p)

		err = TrySetCache(ctx, p, ctx.Response.Body())
		if err != nil {
			doPageErr(pa, ctx, err)
			// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

	})
}

// RegistPost router
func (s *SiteServer) RegistPost(rule string, p IPage) {
	s.Router.POST(rule, func(ctx *fasthttp.RequestCtx) {
		// ctx.SetStatusCode(fasthttp.StatusForbidden)
	})
}

func doPageErr(pd *PageAttr, ctx *fasthttp.RequestCtx, err error) {
	log.App.Error(err)
	var errorURL string
	if len(pd.ErrorURL) > 0 {
		errorURL = pd.ErrorURL

	} else {
		errorURL = pd.SiteData.ErrorURL
	}

	if errorURL == "" {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	// ErrorURL redirect to
	ctx.Redirect(errorURL+"?msg="+err.Error(), fasthttp.StatusSeeOther)
}

// RequireJSWithSrc render
func RequireJSWithSrc(pd *PageAttr, src []byte) *RequireJS {
	return NewRequireJS(pd, "", src)
}

// NewRequireJS render
func NewRequireJS(pd *PageAttr, fileName string, src []byte) *RequireJS {
	attr := make([][2]string, 3, 5)
	attr[0] = [2]string{"defer", ""}
	attr[1] = [2]string{"async", "true"}
	if RunMode == ModePro {
		attr[2] = [2]string{"data-main", "/assets/js/" + RunMode + "/page/" + pd.Dir + "/" + pd.Name + ".min"}
	} else {
		attr[2] = [2]string{"data-main", "/assets/js/" + RunMode + "/page/" + pd.Dir + "/" + pd.Name}
	}
	rj := &RequireJS{}
	rj.Src = src
	rj.Attr = attr
	rj.FileName = fileName

	return rj
}
