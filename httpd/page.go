package httpd

import (
	"errors"
	"time"

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
		if s.Pool == 0 {
			pageHandle(p, ctx)
			return
		}

		done := make(chan struct{})
		if err := pool.Invoke(PoolParams{Typ: invokePage, Page: p, Ctx: ctx, Done: done}); err != nil {
			doAPIError(ctx, errors.New("Throttle limit error"))
		}
		<-done
		select {
		case <-done:
		case <-time.After(s.Timeout):
			ctx.TimeoutError("timeout!")
		}

	})
}

// RegistPost router
func (s *SiteServer) RegistPost(rule string, p IPage) {
	s.Router.POST(rule, func(ctx *fasthttp.RequestCtx) {
		// ctx.SetStatusCode(fasthttp.StatusForbidden)
	})
}

func doPageErr(errorURL string, ctx *fasthttp.RequestCtx, err error) {
	if errorURL == "" {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	// ErrorURL redirect to
	ctx.Redirect(errorURL+"?msg="+err.Error(), fasthttp.StatusSeeOther)
}
