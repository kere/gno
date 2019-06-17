package httpd

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// IPage interface
type IPage interface {
	Data() *PageData
	SetData(*PageData)
	Auth(ctx *fasthttp.RequestCtx) error
	// Before(ctx *fasthttp.RequestCtx) error
	Page(ctx *fasthttp.RequestCtx) error
}

// RegistGet router
func (s *SiteServer) RegistGet(rule string, p IPage) {
	s.Router.GET(rule, func(ctx *fasthttp.RequestCtx) {
		pageHandle(p, ctx)
		// done := make(chan bool)
		// if err := pool.Invoke(PoolParams{Typ: invokePage, Page: p, Ctx: ctx, Done: done}); err != nil {
		// 	doAPIError(ctx, errors.New("Throttle limit error"))
		// }
		// <-done
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
