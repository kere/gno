package httpd

import (
	"net/url"
	"path/filepath"

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
	Data() *PageData
	SetData(*PageData)
	Auth(ctx *fasthttp.RequestCtx) error
	Page(ctx *fasthttp.RequestCtx) error
}

// PageData class
type PageData struct {
	Name, Dir string
	Title     string
	// Secret, Nonce string
	CacheOption PageCacheOption
	JSPosition  int

	SiteData *SiteData
	ErrorURL string
	LoginURL string

	Head []IRender

	JS  []IRenderWith
	CSS []IRenderWith

	Top    []IRender
	Body   []IRender
	Bottom []IRender

	RenderRelease []IRenderRelease
}

// Init page
func (d *PageData) Init(title, name, dir string) {
	d.Title = title
	d.Name = name
	d.Dir = dir
}

// P class
type P struct {
	D PageData
}

// Data page
func (p *P) Data() *PageData {
	if len(p.D.Body) == 0 {
		p.D.Body = append(p.D.Body, &Template{FileName: filepath.Join(p.D.Dir, p.D.Name+DefaultTemplateSubfix)})
	}

	return &p.D
}

// SetData page
func (p *P) SetData(pd *PageData) {
	p.D = *pd
}

// Page do
func (p *P) Page(ctx *fasthttp.RequestCtx) error {
	return nil
}

// Auth page
func (p *P) Auth(ctx *fasthttp.RequestCtx) error {
	return nil
}

// RegistGet router
func (s *SiteServer) RegistGet(rule string, p IPage) {
	s.Router.GET(rule, func(ctx *fasthttp.RequestCtx) {
		pd := p.Data()
		pd.SiteData = s.SiteData
		// do auth
		err := p.Auth(ctx)
		if err != nil {
			var loginURL string
			if len(loginURL) == 0 {
				loginURL = pd.SiteData.LoginURL
			} else {
				loginURL = pd.LoginURL
			}
			u, _ := url.Parse(loginURL)
			u.Query().Add(sAuthURL, string(ctx.RequestURI()))
			ctx.Redirect(u.String(), fasthttp.StatusSeeOther)
			return
		}

		// try cache
		if TryCache(ctx, p) {
			return
		}

		// do page
		err = p.Page(ctx)

		if err != nil {
			doPageErr(pd, ctx, err)
			// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		// fmt.Println("PageName:", pd.Name, reflect.TypeOf(p))

		err = renderPage(ctx, pd, ctx.URI().PathOriginal())
		l := len(pd.RenderRelease)
		if l > 0 {
			for i := 0; i < l; i++ {
				pd.RenderRelease[i].Release()
			}
		}

		if err != nil {
			doPageErr(pd, ctx, err)
			// ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		err = TrySetCache(ctx, p, ctx.Response.Body())
		if err != nil {
			doPageErr(pd, ctx, err)
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

func doPageErr(pd *PageData, ctx *fasthttp.RequestCtx, err error) {
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

// RequireJS render
func RequireJS(pd *PageData, src []byte) *JS {
	requireAttr := make([][2]string, 3, 5)
	requireAttr[0] = [2]string{"defer", ""}
	requireAttr[1] = [2]string{"async", "true"}
	requireAttr[2] = [2]string{"data-main", "/assets/js/" + RunMode + "/page/" + pd.Dir + "/" + pd.Name}

	return &JS{Src: src, Attr: requireAttr}
}
